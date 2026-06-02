package tests

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/app"
	"github.com/canta2899/logo-ls/fs/osfs"
	"github.com/canta2899/logo-ls/model"
)

func TestRealFS_Smoke(t *testing.T) {
	dir := t.TempDir()

	names := []string{"alpha.txt", "beta.go", "gamma.md"}
	for _, n := range names {
		if err := os.WriteFile(filepath.Join(dir, n), []byte("x"), 0o644); err != nil {
			t.Fatalf("write %s: %v", n, err)
		}
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0o755); err != nil {
		t.Fatalf("mkdir sub: %v", err)
	}

	cfg, _, err := app.BuildConfig([]string{"logo-ls", "-1e", dir})
	if err != nil {
		t.Fatalf("BuildConfig: %v", err)
	}

	var stdout, stderr bytes.Buffer
	a := &app.App{
		Config:   cfg,
		Writer:   &stdout,
		Logger:   log.New(&stderr, "", 0),
		FS:       osfs.New(),
		ExitCode: model.CodeOk,
	}
	a.Run()

	if a.ExitCode != model.CodeOk {
		t.Errorf("exit code: want 0, got %d, stderr=%q", a.ExitCode, stderr.String())
	}
	if s := strings.TrimSpace(stderr.String()); s != "" {
		t.Errorf("expected silent stderr, got: %q", s)
	}
	out := stripANSI(stdout.String())
	for _, n := range append(names, "sub") {
		if !strings.Contains(out, n) {
			t.Errorf("expected %q in output:\n%s", n, out)
		}
	}
}
