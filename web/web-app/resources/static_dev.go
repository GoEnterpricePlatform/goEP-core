//go:build dev

package resources

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"
)

var moduleStaticName = regexp.MustCompile(`^[a-z][a-z0-9]*(-[a-z0-9]+)*$`)

func Handler() http.Handler {
	slog.Info("static assets are being served directly", "path", StaticDirectoryPath)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		requestPath := strings.TrimPrefix(r.URL.Path, "/static/")
		cleanPath := path.Clean(requestPath)
		if cleanPath != "." && cleanPath != requestPath {
			http.NotFound(w, r)
			return
		}
		if strings.HasPrefix(cleanPath, "modules/") {
			parts := strings.SplitN(strings.TrimPrefix(cleanPath, "modules/"), "/", 2)
			if len(parts) != 2 || !moduleStaticName.MatchString(parts[0]) || parts[1] == "module-assets.json" {
				http.NotFound(w, r)
				return
			}
			moduleDirectory := "addons/" + parts[0] + "/static/dist"
			if info, err := os.Stat(moduleDirectory); err != nil || !info.IsDir() {
				http.NotFound(w, r)
				return
			}
			prefix := fmt.Sprintf("/static/modules/%s/", parts[0])
			http.StripPrefix(prefix, http.FileServerFS(os.DirFS(moduleDirectory))).ServeHTTP(w, r)
			return
		}
		http.StripPrefix("/static/", http.FileServerFS(os.DirFS(StaticDirectoryPath))).ServeHTTP(w, r)
	})
}

func StaticPath(path string) string {
	return "/static/" + path
}
