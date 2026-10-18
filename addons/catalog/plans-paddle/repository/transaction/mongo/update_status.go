package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans-paddle/model/transaction"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *TransactionRepository) UpsertFromWebhook(ctx context.Context, transaction *domain.PaddleTransaction, eventAt time.Time) error {
	if transaction == nil || transaction.PaddleTransactionID == "" {
		return fmt.Errorf("Paddle webhook transaction is missing Paddle transaction ID")
	}
	if transaction.ID == "" {
		transaction.ID = bson.NewObjectID().Hex()
	}
	transaction.LastEventAt = eventAt
	doc := model.FromDomainPaddleTransaction(transaction)
	var existing model.PaddleTransactionModel
	findErr := r.Collection.FindOne(ctx, bson.M{"paddle_transaction_id": doc.PaddleTransactionID}).Decode(&existing)
	if findErr == nil && existing.PaymentDetailsProvided {
		doc.PaymentDetailsProvided = true
		doc.CurrencyCode = existing.CurrencyCode
		doc.Subtotal = existing.Subtotal
		doc.Tax = existing.Tax
		doc.Total = existing.Total
	} else if findErr != nil && findErr != mongo.ErrNoDocuments {
		return fmt.Errorf("error loading Paddle transaction payment snapshot: %w", findErr)
	}
	set := bson.M{
		"status":        doc.Status,
		"updated_at":    doc.UpdatedAt,
		"last_event_at": eventAt,
	}
	if doc.UserID != "" {
		set["user_id"] = doc.UserID
	}
	if doc.CustomerID != "" {
		set["customer_id"] = doc.CustomerID
	}
	if doc.PlanID != "" {
		set["plan_id"] = doc.PlanID
	}
	if doc.PlanItemID != "" {
		set["plan_item_id"] = doc.PlanItemID
	}
	if doc.SubscriptionID != nil {
		set["subscription_id"] = *doc.SubscriptionID
	}
	// Keep the transaction's recorded payment amounts from Paddle. These values
	// are snapshots of the payment and must never be recalculated from catalog data.
	if doc.PaymentDetailsProvided {
		set["payment_details_provided"] = true
		set["currency_code"] = doc.CurrencyCode
		set["subtotal"] = doc.Subtotal
		set["tax"] = doc.Tax
		set["total"] = doc.Total
	}

	filter := bson.M{
		"paddle_transaction_id": doc.PaddleTransactionID,
		"$or": bson.A{
			bson.M{"last_event_at": bson.M{"$lte": eventAt}},
			bson.M{"last_event_at": bson.M{"$exists": false}},
		},
	}
	update := bson.M{
		"$set":         set,
		"$setOnInsert": bson.M{"_id": doc.ID, "paddle_transaction_id": doc.PaddleTransactionID},
	}
	result, err := r.Collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			// A newer event or concurrent delivery may have inserted this Paddle ID.
			_, retryErr := r.Collection.UpdateOne(ctx, filter, update)
			if retryErr == nil {
				return nil
			}
		}
		return fmt.Errorf("error upserting Paddle transaction from webhook: %w", err)
	}
	_ = result
	return nil
}
