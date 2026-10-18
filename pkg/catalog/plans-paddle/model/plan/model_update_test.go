package model

import (
	"testing"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func TestFromDomainPaddlePlanKeepsExistingAndCreatesNewItemIDs(t *testing.T) {
	planID := bson.NewObjectID()
	existingID := bson.NewObjectID()
	input := &domain.PaddlePlan{PlanID: planID.Hex(), Items: []*domain.PaddlePlanItem{
		{ID: existingID.Hex(), PlanItemID: bson.NewObjectID().Hex()},
		{PlanItemID: bson.NewObjectID().Hex()},
	}}
	m, err := FromDomainPaddlePlan(input, bson.NewObjectID())
	if err != nil {
		t.Fatal(err)
	}
	if m.Items[0].ID != existingID || m.Items[1].ID.IsZero() || m.Items[1].ID == existingID {
		t.Fatal("existing item ID was lost or new item ID was not generated")
	}
}
