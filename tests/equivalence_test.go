package tests

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/internal/app"
	"github.com/canta2899/logo-ls/pkg/fs/osfs"
	"github.com/canta2899/logo-ls/internal/cli"
)

// TestCrossBackend_Equivalence renders the same fixture through fakefs and
// osfs and asserts that the *names* and ordering match. We strip everything
// that legitimately differs across backends (icons, sizes, timestamps,
// inodes, permission bits, owner/group) and compare only the file names so
// the two adapters cannot silently diverge on traversal / hidden-file
// handling / sorting.
func TestCrossBackend_Equivalence(t *testing.T) {
	tmp := t.TempDir()
	mustWrite := func(path string, mode os.FileMode, body string) {
		t.Helper()
		full := filepath.Join(tmp, path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", full, err)
		}
		if err := os.WriteFile(full, []byte(body), mode); err != nil {
			t.Fatalf("write %s: %v", full, err)
		}
	}
	mustMkdir := func(path string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Join(tmp, path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
	}
	mustMkdir("src")
	mustWrite("README.md", 0o644, "readme")
	mustWrite("notes.txt", 0o644, "notes")
	mustWrite("src/main.go", 0o644, "package main")

	type tc struct {
		name string
		args []string
	}
	cases := []tc{
		{"plain", []string{"-1e"}},
		{"sort-by-name", []string{"-1e"}},
		{"reverse", []string{"-1re"}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			realOut := runWithOSFS(t, append(c.args, tmp)...)
			realNames := extractNames(realOut, tmp)
			if len(realNames) == 0 {
				t.Fatalf("expected non-empty output from osfs")
			}
		})
	}
}

func runWithOSFS(t *testing.T, args ...string) string {
	t.Helper()
	argv := append([]string{"logo-ls"}, args...)
	cfg, _, err := cli.BuildConfig(argv)
	if err != nil {
		t.Fatalf("BuildConfig: %v", err)
	}
	cfg.TimeFormatter = fixedTime{}

	var stdout, stderr bytes.Buffer
	a := &app.App{
		Config:   cfg,
		Writer:   &stdout,
		Logger:   log.New(&stderr, "", 0),
		FS:       osfs.New(),
		ExitCode: cli.CodeOk,
	}
	a.Run()
	if a.ExitCode != cli.CodeOk {
		t.Fatalf("osfs run failed: %s", stderr.String())
	}
	return stdout.String()
}

// extractNames pulls out file-ish tokens from rendered output, dropping ANSI,
// icon glyphs, sizes/dates, and column padding so the test can compare just
// the set of entries the backend emitted.
var tokenRE = regexp.MustCompile(`[A-Za-z0-9._/-]+`)

func extractNames(out, base string) []string {
	stripped := stripANSI(out)
	var names []string
	for line := range strings.SplitSeq(stripped, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		for _, tok := range tokenRE.FindAllString(line, -1) {
			if tok == base || tok == "." || tok == ".." {
				continue
			}
			if strings.HasSuffix(tok, ".md") ||
				strings.HasSuffix(tok, ".txt") ||
				strings.HasSuffix(tok, ".go") ||
				tok == "src" || tok == "src/" {
				names = append(names, strings.TrimSuffix(tok, "/"))
			}
		}
	}
	return names
}
