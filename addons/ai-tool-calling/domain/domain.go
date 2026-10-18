package domain

type ToolSelectionResult struct {
	Operation    string
	Action       string
	Arguments    map[string]any
	Fields       []FieldInfo
	Required     []string
	NeedsConfirm bool
	Message      string
}

type FieldInfo struct {
	Name     string
	Value    string
	Required bool
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
