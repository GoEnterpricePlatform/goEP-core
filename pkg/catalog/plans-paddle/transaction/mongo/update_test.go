package mongo

import (
	"testing"

	paddleD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	planD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
)

func TestLinkUpdatedItems(t *testing.T) {
	plan := &planD.Plan{Items: []*planD.PlanItem{{ID: "existing"}, {ID: "generated"}}}
	paddle := &paddleD.PaddlePlan{Items: []*paddleD.PaddlePlanItem{{PlanItemID: "old"}, {}}}
	if err := linkUpdatedItems(plan, paddle); err != nil {
		t.Fatal(err)
	}
	if paddle.Items[0].PlanItemID != "existing" || paddle.Items[1].PlanItemID != "generated" {
		t.Fatal("paddle items were not linked to updated plan items")
	}
}
