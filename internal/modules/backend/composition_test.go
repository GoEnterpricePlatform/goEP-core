package backend

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestModuleRoutesExposeOnlyEnabledModules(t *testing.T) {
	apiMux := http.NewServeMux()

	enabledMux := moduleRoutes(apiMux, map[string]bool{"identity": true}, "identity")
	if enabledMux != apiMux {
		t.Fatal("enabled module should register routes on the public API mux")
	}

	disabledMux := moduleRoutes(apiMux, map[string]bool{"identity": true}, "posts")
	if disabledMux == apiMux {
		t.Fatal("disabled module should receive an isolated mux")
	}
	disabledMux.HandleFunc("GET /posts", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	response := httptest.NewRecorder()
	apiMux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/posts", nil))
	if response.Code != http.StatusNotFound {
		t.Fatalf("public API status = %d, want %d for disabled module", response.Code, http.StatusNotFound)
	}

	response = httptest.NewRecorder()
	disabledMux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/posts", nil))
	if response.Code != http.StatusNoContent {
		t.Fatalf("isolated module status = %d, want %d", response.Code, http.StatusNoContent)
	}
}
