package core

import "testing"

func TestUpdateAllowsNewItemsAndRejectsDuplicateCombinations(t *testing.T) {
	req := UpdatePaddlePlanReq{
		Name: "Plan", PaddleProductID: "pro_1",
		Items: []*UpdatePaddlePlanItem{
			{ID: "existing", Features: []string{"A"}, VarOptionIDs: []string{"a"}, PaddlePriceID: "pri_1"},
			{Features: []string{"B"}, VarOptionIDs: []string{"b"}, PaddlePriceID: "pri_2"},
		},
	}
	if err := req.Validate(); err != nil {
		t.Fatal(err)
	}
	req.Items[1].VarOptionIDs = []string{"a"}
	if err := req.Validate(); err == nil {
		t.Fatal("duplicate combination should be rejected")
	}
}
