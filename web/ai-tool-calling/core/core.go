package core

type ChatMessage struct {
	Role    string
	Content string

	// acciones (create/update/delete)
	//Arguments map[string]any
	//Missing   []string

	// resultados tipo tabla (get_all)
	Table *TableView
}

type TableView struct {
	Columns []string
	Rows    []map[string]any
}

type AssistantPageData struct {
	Messages  []ChatMessage
	Prompt    string
	Error     string
	Arguments map[string]any
	Missing   []string
	/* 	ErrorSubTmpl string
	Success      string */
}

/* type ConversationContext struct {
	Prompt    string
	Operation string
	Data      map[string]string
	Missing   []string
	Message   string
	Completed bool
} */
