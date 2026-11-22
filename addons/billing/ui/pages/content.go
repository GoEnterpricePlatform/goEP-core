package pages

import "github.com/a-h/templ"

// BillingDashboardContent selects the page associated with the active menu item.
func BillingDashboardContent(activePage string) templ.Component {
	switch activePage {
	case "example-group":
		return BillingExampleGroupPage()
	case "example":
		return BillingExamplePage()
	default:
		return BillingOverviewPage()
	}
}
