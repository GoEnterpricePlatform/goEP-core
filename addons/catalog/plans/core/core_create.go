package core

type CreateOrUpdatePlanReq struct {
	Name        string               `json:"name"`
	Description *string              `json:"description"`
	PlanItems   []*CreatePlanItemReq `json:"plan_items"`
}

type CreatePlanItemReq struct {
	ID           string   `json:"id"`
	Features     []string `json:"features"`
	VarOptionIDs []string `json:"var_option_ids,omitempty"`
}

func (p CreateOrUpdatePlanReq) Validate() error {
	return validatePlanFields(p)
}


func (p CreateOrUpdatePlanReq) ValidateUpdate() error {
	if err := validatePlanFields(p); err != nil {
		return err
	}

	return validatePlanItemIDs(p)
}