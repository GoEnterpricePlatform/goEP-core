package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *SubscriptionRepository) Upsert(ctx context.Context, subscription *domain.PaddleSubscription, lastEventAt time.Time) error {
	if subscription == nil || subscription.PaddleSubscriptionID == "" {
		return fmt.Errorf("Paddle subscription is missing Paddle subscription ID")
	}
	if subscription.ID == "" {
		var existing model.PaddleSubscriptionModel
		if err := r.Collection.FindOne(ctx, bson.M{"paddle_subscription_id": subscription.PaddleSubscriptionID}).Decode(&existing); err == nil {
			subscription.ID = existing.ID
		} else if err != mongo.ErrNoDocuments {
			return fmt.Errorf("error loading Paddle subscription before upsert: %w", err)
		}
	}
	if subscription.ID == "" {
		subscription.ID = bson.NewObjectID().Hex()
	}
	doc := model.FromDomainPaddleSubscription(subscription, lastEventAt)
	set := bson.M{
		"status":        doc.Status,
		"updated_at":    doc.UpdatedAt,
		"last_event_at": doc.LastEventAt,
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
	setOptionalString(set, "current_period_starts_at", doc.CurrentPeriodStartsAt)
	setOptionalString(set, "current_period_ends_at", doc.CurrentPeriodEndsAt)
	setOptionalString(set, "next_billed_at", doc.NextBilledAt)
	setOptionalString(set, "canceled_at", doc.CanceledAt)
	if doc.ScheduledChange != nil {
		set["scheduled_change"] = doc.ScheduledChange
	} else {
		set["scheduled_change"] = nil
	}
	filter := bson.M{
		"paddle_subscription_id": doc.PaddleSubscriptionID,
		"$or": bson.A{
			bson.M{"last_event_at": bson.M{"$lte": lastEventAt}},
			bson.M{"last_event_at": bson.M{"$exists": false}},
		},
	}
	update := bson.M{
		"$set":         set,
		"$setOnInsert": bson.M{"_id": doc.ID, "paddle_subscription_id": doc.PaddleSubscriptionID},
	}
	_, err := r.Collection.UpdateOne(ctx, filter, update, options.UpdateOne().SetUpsert(true))
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			if _, retryErr := r.Collection.UpdateOne(ctx, filter, update); retryErr == nil {
				return nil
			}
		}
		return fmt.Errorf("error upserting Paddle subscription: %w", err)
	}
	return nil
}
