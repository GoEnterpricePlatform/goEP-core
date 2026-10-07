package service

import (
	"context"
	"errors"
	"testing"

	variationD "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

type variationRepoStub struct{ deleted bool }

func (*variationRepoStub) Insert(context.Context, *variationD.Variation) error { return nil }
func (*variationRepoStub) FindAllWithOptions(context.Context) ([]*variationD.Variation, error) {
	return nil, nil
}
func (*variationRepoStub) Update(context.Context, *variationD.Variation) error { return nil }
func (r *variationRepoStub) Delete(context.Context, string) error              { r.deleted = true; return nil }

type optionRepoStub struct{ hasOptions, deleted bool }

func (*optionRepoStub) Insert(context.Context, *variationD.VarOption) error { return nil }
func (*optionRepoStub) Update(context.Context, *variationD.VarOption) error { return nil }
func (r *optionRepoStub) Delete(context.Context, string, string) error      { r.deleted = true; return nil }
func (*optionRepoStub) FindByIDs(context.Context, []string) ([]*variationD.VarOption, error) {
	return nil, nil
}
func (r *optionRepoStub) HasByVariation(context.Context, string) (bool, error) {
	return r.hasOptions, nil
}

type optionUsageStub struct{ inUse bool }

func (r *optionUsageStub) IsOptionInUse(context.Context, string) (bool, error) {
	return r.inUse, nil
}

func assertConflict(t *testing.T, err error) {
	t.Helper()
	var appErr sharedD.AppError
	if !errors.As(err, &appErr) || appErr.Code != sharedD.ErrCodeConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
}

func TestDeleteVariationBlocksOptions(t *testing.T) {
	variationRepo := &variationRepoStub{}
	optionRepo := &optionRepoStub{hasOptions: true}
	srv := NewVariationSrv(variationRepo, optionRepo, &optionUsageStub{})
	assertConflict(t, srv.DeleteVariation(context.Background(), "variation"))
	if variationRepo.deleted {
		t.Fatal("variation with options was deleted")
	}
	optionRepo.hasOptions = false
	if err := srv.DeleteVariation(context.Background(), "variation"); err != nil {
		t.Fatal(err)
	}
	if !variationRepo.deleted {
		t.Fatal("variation without options was not deleted")
	}
}

func TestDeleteOptionBlocksPlanReference(t *testing.T) {
	optionRepo := &optionRepoStub{}
	usage := &optionUsageStub{inUse: true}
	srv := NewVariationSrv(&variationRepoStub{}, optionRepo, usage)
	assertConflict(t, srv.DeleteVarOption(context.Background(), "option", "variation"))
	if optionRepo.deleted {
		t.Fatal("referenced option was deleted")
	}
	usage.inUse = false
	if err := srv.DeleteVarOption(context.Background(), "option", "variation"); err != nil {
		t.Fatal(err)
	}
	if !optionRepo.deleted {
		t.Fatal("unused option was not deleted")
	}
}
