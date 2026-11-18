package layouts

import (
	"context"
	"io"

	"github.com/GoEnterpricePlatform/goEP-core/internal/modules/navigation"
	"github.com/a-h/templ"
)

const goEPModuleHeader = `<header class="border-b border-zinc-200"><div class="mx-auto flex max-w-7xl items-center justify-between px-6 py-5 lg:px-8"><a href="/v1/goep-admin/modules" class="text-sm font-bold tracking-[0.2em] text-black">GOEP CORE</a></div></header>`

// GoEPModuleLayout renders a standalone module page with the GoEP Core header.
func GoEPModuleLayout(content templ.Component) templ.Component {
	return BaseLayout("GoEP Core", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="min-h-screen bg-white">`+goEPModuleHeader+`<main class="mx-auto max-w-7xl px-6 py-12 lg:px-8">`); err != nil {
			return err
		}
		if err := content.Render(ctx, w); err != nil {
			return err
		}
		_, err := io.WriteString(w, `</main></div>`)
		return err
	}))
}

// GoEPDashboardLayout adds the GoEP Core header above the dashboard sidebar.
func GoEPDashboardLayout(activePage string, menu []navigation.Item, content templ.Component) templ.Component {
	return BaseLayout("GoEP Core", templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(w, `<div class="min-h-screen bg-white">`+goEPModuleHeader); err != nil {
			return err
		}
		if err := DashboardLayoutContent(activePage, []string{"ROLE_SYSTEM_ADMIN"}, nil, menu, content).Render(ctx, w); err != nil {
			return err
		}
		_, err := io.WriteString(w, `</div>`)
		return err
	}))
}
