package domain

type ToolSelectionResult struct {
	Operation    string         `json:"operation"`
	Arguments    map[string]any `json:"arguments"`
	Fields       []FieldInfo    `json:"fields"`
	Missing      []string       `json:"missing"`
	Optional     []string       `json:"optinal"`
	NeedsConfirm bool           `json:"needs_confirm"`
}

type FieldInfo struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}
