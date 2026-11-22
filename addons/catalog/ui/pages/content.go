package pages

import (
	"strings"

	"github.com/a-h/templ"
)

func DashboardContent(activePage string) templ.Component {
	switch activePage {
	case "variations":
		return CatalogVariationsContent(nil, "")
	case "paddle-plans":
		return CatalogPlaceholderContent("Paddle plans")
	case "polar-plans":
		return CatalogPlaceholderContent("Polar plans")
	default:
		return CatalogGeneralContent()
	}
}

func SectionTitle(activePage string) string {
	label := strings.ReplaceAll(activePage, "-", " ")
	if label == "" {
		return "General"
	}
	return strings.ToUpper(label[:1]) + label[1:]
}
