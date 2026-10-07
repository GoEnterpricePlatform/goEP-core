package mongo

import (
	"context"
	"testing"

	paddleD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	paddleP "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/port"
	planP "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/port"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

type deletePaddleRepo struct {
	paddleP.PaddlePlanRepo
	steps   *[]string
	missing bool
}

func (r *deletePaddleRepo) Find(_ context.Context, id string) (*paddleD.PaddlePlan, error) {
	*r.steps = append(*r.steps, "find:"+id)
	if r.missing {
		return nil, sharedD.ErrNotFound
	}
	return &paddleD.PaddlePlan{PlanID: "linked-plan"}, nil
}
func (r *deletePaddleRepo) Delete(_ context.Context, id string) error {
	*r.steps = append(*r.steps, "paddle:"+id)
	return nil
}

type deleteBasePlanRepo struct {
	planP.PlanRepo
	steps *[]string
}

func (r *deleteBasePlanRepo) Delete(_ context.Context, id string) error {
	*r.steps = append(*r.steps, "plan:"+id)
	return nil
}

func TestDeleteLinkedRemovesBothRecords(t *testing.T) {
	steps := []string{}
	tx := &Transaction{PaddlePlanRepo: &deletePaddleRepo{steps: &steps}, PlanRepo: &deleteBasePlanRepo{steps: &steps}}
	if err := tx.deleteLinked(context.Background(), "paddle-plan"); err != nil {
		t.Fatal(err)
	}
	if len(steps) != 3 || steps[0] != "find:paddle-plan" || steps[1] != "paddle:paddle-plan" || steps[2] != "plan:linked-plan" {
		t.Fatalf("unexpected deletion sequence: %v", steps)
	}
}

func TestDeleteLinkedStopsWhenPaddlePlanIsMissing(t *testing.T) {
	steps := []string{}
	tx := &Transaction{PaddlePlanRepo: &deletePaddleRepo{steps: &steps, missing: true}, PlanRepo: &deleteBasePlanRepo{steps: &steps}}
	if err := tx.deleteLinked(context.Background(), "missing"); err == nil {
		t.Fatal("expected not found error")
	}
	if len(steps) != 1 {
		t.Fatalf("unexpected deletion after missing record: %v", steps)
	}
}
