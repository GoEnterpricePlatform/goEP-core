//go:build !dev

package resources

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"

	backendModules "github.com/GoEnterpricePlatform/goEP-core/internal/modules/backend"
	"github.com/benbjohnson/hashfs"
)

var (
	//go:embed static
	StaticDirectory embed.FS
	StaticSys       = hashfs.NewFS(moduleStaticFS{base: StaticDirectory})
)

// moduleStaticFS maps public module URLs to each addon's embedded static/dist.
// Shared Web App assets continue to come from resources/static.
type moduleStaticFS struct {
	base fs.FS
}

func (f moduleStaticFS) Open(name string) (fs.File, error) {
	const prefix = "static/modules/"
	if !strings.HasPrefix(name, prefix) {
		return f.base.Open(name)
	}

	remainder := strings.TrimPrefix(name, prefix)
	parts := strings.SplitN(remainder, "/", 2)
	if len(parts) != 2 || parts[0] == "" || path.Clean(name) != name {
		return nil, fs.ErrNotExist
	}
	if parts[1] == "module-assets.json" {
		return nil, fs.ErrNotExist
	}
	moduleAssets, ok := backendModules.StaticDirectories[parts[0]]
	if !ok {
		return nil, fs.ErrNotExist
	}
	moduleRoot, err := fs.Sub(moduleAssets, "static/dist")
	if err != nil {
		return nil, fs.ErrNotExist
	}
	return moduleRoot.Open(parts[1])
}

func Handler() http.Handler {
	slog.Debug("static assets are embedded")
	return hashfs.FileServer(StaticSys)
}

func StaticPath(path string) string {
	return "/" + StaticSys.HashName("static/"+path)
}
