package domain

// ! ver con que nos vamos a quedar
type ToolSelectionResult struct {
	Operation    string         `json:"operation"`
	Action       string         `json:"action"`
	Arguments    map[string]any `json:"arguments"`
	Fields       []FieldInfo    `json:"fields"`
	Required     []string       `json:"required"`
	NeedsConfirm bool           `json:"needs_confirm"`
}

type FieldInfo struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Required bool   `json:"required"`
}

type ChatMessage struct {
	Role    string
	Content string
	Table   *TableView
}

type TableView struct {
	Columns []string
	Rows    []map[string]any
}
