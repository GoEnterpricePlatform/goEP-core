package service

import (
	"context"
	"testing"
	"time"

	paddleD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans-paddle/domain"
	planD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
)

type combinationPaddleRepo struct{ current *paddleD.PaddlePlan }

func (r *combinationPaddleRepo) Insert(context.Context, *paddleD.PaddlePlan) error { return nil }
func (r *combinationPaddleRepo) Find(context.Context, string) (*paddleD.PaddlePlan, error) {
	return r.current, nil
}
func (r *combinationPaddleRepo) FindAll(context.Context, int64, int64) ([]*paddleD.PaddlePlan, error) {
	return nil, nil
}
func (r *combinationPaddleRepo) Count(context.Context) (int64, error) { return 0, nil }
func (r *combinationPaddleRepo) Update(context.Context, string, *paddleD.PaddlePlan) error {
	return nil
}
func (r *combinationPaddleRepo) Delete(context.Context, string) error { return nil }

type combinationPlanRepo struct{ current *planD.Plan }

func (r *combinationPlanRepo) Insert(context.Context, *planD.Plan) error { return nil }
func (r *combinationPlanRepo) Find(context.Context, string) (*planD.Plan, error) {
	return r.current, nil
}
func (r *combinationPlanRepo) FindAll(context.Context, int64, int64) ([]*planD.Plan, error) {
	return nil, nil
}
func (r *combinationPlanRepo) Count(context.Context) (int64, error)              { return 0, nil }
func (r *combinationPlanRepo) Update(context.Context, string, *planD.Plan) error { return nil }
func (r *combinationPlanRepo) Delete(context.Context, string) error              { return nil }

type combinationTx struct {
	plan   *planD.Plan
	paddle *paddleD.PaddlePlan
}

func (t *combinationTx) Create(context.Context, *planD.Plan, *paddleD.PaddlePlan) error { return nil }
func (t *combinationTx) Update(_ context.Context, plan *planD.Plan, paddle *paddleD.PaddlePlan) error {
	t.plan, t.paddle = plan, paddle
	return nil
}
func (t *combinationTx) Delete(context.Context, string) error { return nil }

func TestUpdateReplacesCombinations(t *testing.T) {
	created := time.Now().UTC()
	tx := &combinationTx{}
	s := &Service{
		PaddlePlanTx: tx,
		PaddlePlanRepo: &combinationPaddleRepo{current: &paddleD.PaddlePlan{ID: "paddle", PlanID: "plan", Items: []*paddleD.PaddlePlanItem{
			{ID: "kept", PlanItemID: "base-kept", CreatedAt: &created},
			{ID: "removed", PlanItemID: "base-removed", CreatedAt: &created},
		}}},
		PlanRepo: &combinationPlanRepo{current: &planD.Plan{ID: "plan", Items: []*planD.PlanItem{
			{ID: "base-kept", Status: planD.PlanStatusInactive, CreatedAt: &created},
			{ID: "base-removed", Status: planD.PlanStatusActive, CreatedAt: &created},
		}}},
	}
	input := &paddleD.PaddlePlan{Items: []*paddleD.PaddlePlanItem{
		{ID: "kept", PaddlePriceID: "pri_kept"},
		{PaddlePriceID: "pri_new"},
	}}
	if err := s.Update(context.Background(), "paddle", input); err != nil {
		t.Fatal(err)
	}
	if len(tx.plan.Items) != 2 || len(tx.paddle.Items) != 2 {
		t.Fatal("item set was not replaced")
	}
	if tx.plan.Items[0].ID != "base-kept" || tx.plan.Items[0].Status != planD.PlanStatusInactive || tx.paddle.Items[0].ID != "kept" {
		t.Fatal("existing combination was not preserved")
	}
	if tx.plan.Items[1].ID != "" || tx.plan.Items[1].Status != planD.PlanStatusActive || tx.paddle.Items[1].ID != "" || tx.paddle.Items[1].CreatedAt == nil {
		t.Fatal("new combination was not prepared for insertion")
	}
}
