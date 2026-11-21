package components

import "github.com/GoEnterpricePlatform/goEP-core/internal/modules/navigation"

func hasActiveChild(activePage string, children []navigation.Item) bool {
	for _, child := range children {
		if child.ActivePage == activePage || hasActiveChild(activePage, child.Children) {
			return true
		}
	}
	return false
}
