package mongo

import (
	"context"
	"fmt"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/model"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) Update(ctx context.Context, id string, paddlePlan *domain.PaddlePlan) error {
	oID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return sharedD.ErrIncorrectID
	}
	m, err := model.FromDomainPaddlePlan(paddlePlan, oID)
	if err != nil {
		return err
	}
	result, err := r.Collection.UpdateOne(ctx, bson.M{"_id": oID, "plan_id": m.PlanID}, bson.M{"$set": bson.M{
		"paddle_product_id": m.PaddleProductID,
		"order":             m.Order,
		"items":             m.Items,
		"updated_at":        m.UpdatedAt,
	}})
	if err != nil {
		return fmt.Errorf("error updating paddle plan: %w", err)
	}
	if result.MatchedCount == 0 {
		return sharedD.ErrNotFound
	}
	for i, item := range m.Items {
		paddlePlan.Items[i].ID = item.ID.Hex()
	}
	return nil
}
