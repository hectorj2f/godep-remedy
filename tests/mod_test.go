package tests

import (
	"os/exec"
	"path"
	"testing"

	"github.com/hectorj2f/godep-remedy/cmd"
	"github.com/hectorj2f/godep-remedy/pkg/types"
)

func TestRun(t *testing.T) {
	fulcioModules := make(map[string]*types.Module, 0)
	fulcioModules["github.com/go-jose/go-jose/v3"] = &types.Module{
		Name:    "github.com/go-jose/go-jose/v3",
		Version: "v3.0.1",
	}
	calicoModules := make(map[string]*types.Module, 0)
	calicoModules["k8s.io/kubernetes"] = &types.Module{
		Name:    "k8s.io/kubernetes",
		Version: "v1.26.11",
	}

	caModules := make(map[string]*types.Module, 0)
	caModules["go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"] = &types.Module{
		Name:    "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
		Version: "v0.46.1",
	}
	caModules["go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"] = &types.Module{
		Name:    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
		Version: "v0.46.1",
	}
	caModules["go.opentelemetry.io/otel/sdk"] = &types.Module{
		Name:    "go.opentelemetry.io/otel/sdk",
		Version: "v1.21.0",
	}
	caModules["go.opentelemetry.io/otel"] = &types.Module{
		Name:    "go.opentelemetry.io/otel",
		Version: "v1.21.0",
	}
	caModules["go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"] = &types.Module{
		Name:    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
		Version: "v1.21.0",
	}
	caModules["k8s.io/kubernetes"] = &types.Module{
		Name:    "k8s.io/kubernetes",
		Version: "v1.25.16",
	}

	testCases := []struct {
		name        string
		modules     map[string]*types.Module
		project     string
		modRootPath string
	}{
		{
			name:        "no downgrade",
			modules:     fulcioModules,
			project:     "fulcio",
			modRootPath: "./fulcio",
		},
		{
			name:        "no downgrade",
			modules:     calicoModules,
			project:     "calico",
			modRootPath: "./calico",
		},
		{
			name:        "no downgrade",
			modules:     caModules,
			project:     "cluster-autoscaler",
			modRootPath: "../../autoscaler/cluster-autoscaler",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpdir := t.TempDir()
			copyFile(t, path.Join(tc.modRootPath, "go.mod"), tmpdir)
			//copyFile(t, path.Join(tc.modRootPath, "go.sum"), tmpdir)
			//t.Logf("tmpdir: %s", tmpdir)
			err := cmd.Run(tc.modRootPath, tc.modules, true, true, false)
			if err != nil {
				t.Fatalf("found one error for project %s with message: %v", tc.project, err)
			}
			//t.Fatal(err)
		})
	}
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	_, err := exec.Command("cp", "-r", src, dst).Output()
	if err != nil {
		t.Fatal(err)
	}
}
