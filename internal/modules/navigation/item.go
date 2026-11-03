package navigation

// Item describes one module-owned entry in a shared admin sidebar.
type Item struct {
	Title      string
	Href       string
	ActivePage string
	Permission string
}
