package modules

import (
	"strings"
	"testing"
)

type runtimeTestContract interface {
	Value() string
}

type runtimeTestValue string

func (v runtimeTestValue) Value() string { return string(v) }

func TestServicesProvideAndRequireTypedContract(t *testing.T) {
	services := NewServices()
	if err := ProvideService[runtimeTestContract](services, runtimeTestValue("available")); err != nil {
		t.Fatalf("ProvideService() error = %v", err)
	}
	got, err := RequireService[runtimeTestContract](services)
	if err != nil {
		t.Fatalf("RequireService() error = %v", err)
	}
	if got.Value() != "available" {
		t.Fatalf("RequireService() value = %q, want %q", got.Value(), "available")
	}
	if err := ProvideService[runtimeTestContract](services, runtimeTestValue("duplicate")); err == nil || !strings.Contains(err.Error(), "already provided") {
		t.Fatalf("ProvideService(duplicate) error = %v, want duplicate-service error", err)
	}
}

func TestRequireServiceReportsMissingContract(t *testing.T) {
	_, err := RequireService[runtimeTestContract](NewServices())
	if err == nil || !strings.Contains(err.Error(), "is not available") {
		t.Fatalf("RequireService(missing) error = %v, want missing-service error", err)
	}
}

func TestValidateBackendRequiresRuntimeInitializer(t *testing.T) {
	registry := NewRegistry()
	if err := RegisterFactory(registry, Factory[string, string]{
		Name: "reports",
		New:  func(value string) (string, error) { return value, nil },
	}); err != nil {
		t.Fatalf("RegisterFactory() error = %v", err)
	}
	manifest := []Manifest{{Name: "reports", Version: "1.0.0"}}
	if err := registry.Validate(manifest); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
	if err := registry.ValidateBackend(manifest); err == nil || !strings.Contains(err.Error(), "no runtime initializer") {
		t.Fatalf("ValidateBackend() error = %v, want missing-initializer error", err)
	}
	if err := RegisterBackendInitializer(registry, "reports", func(BackendRuntime) (any, error) {
		return "reports", nil
	}); err != nil {
		t.Fatalf("RegisterBackendInitializer() error = %v", err)
	}
	if err := registry.ValidateBackend(manifest); err != nil {
		t.Fatalf("ValidateBackend() with initializer error = %v", err)
	}
}
