// Package tests contains the end-to-end test suite for logo-ls.
package tests

import (
	"bytes"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/canta2899/logo-ls/internal/app"
	"github.com/canta2899/logo-ls/internal/cli"
	"github.com/canta2899/logo-ls/pkg/fs"
)

var updateGolden = flag.Bool("update", false, "regenerate golden files")

// TestMain pins the collation locale to C for the whole suite so name sorting
// is deterministic regardless of the host's locale. logo-ls sorts names via
// the active LC_COLLATE (matching system ls), so without this the goldens
// would drift between C and UTF-8 environments.
func TestMain(m *testing.M) {
	os.Setenv("LC_ALL", "C")
	os.Setenv("LC_COLLATE", "C")
	os.Setenv("LANG", "C")
	os.Exit(m.Run())
}

// fixedTime is a deterministic time formatter. It ignores time.Now() so
// goldens don't drift over the years.
type fixedTime struct{}

func (fixedTime) Format(t *time.Time) string {
	return t.Format("2006-01-02 15:04")
}

// runResult is the outcome of a single logo-ls invocation.
type runResult struct {
	Stdout   string
	Stderr   string
	ExitCode cli.ExitCode
}

// runApp builds an App with the given FS, parses the given flag arguments,
// runs the app, and returns stdout/stderr/exit. It never calls os.Exit.
//
// args should NOT include argv[0]; the harness prepends a synthetic program
// name. Positional path arguments default to "/root" when omitted so tests
// don't need to repeat the fixture root.
func runApp(t *testing.T, vfs fs.FS, args ...string) runResult {
	t.Helper()

	// Force --no-override so the harness never picks up the developer's
	// personal icon overrides from $HOME. Keeps goldens reproducible.
	argv := append([]string{"logo-ls", "--no-override"}, args...)
	cfg, _, err := cli.BuildConfig(argv)
	if err != nil {
		t.Fatalf("BuildConfig(%v): %v", args, err)
	}
	cfg.TimeFormatter = fixedTime{}

	// If the user did not pass a path argument, target "/root" (the
	// conventional fixture root used by tests in this package).
	if len(cfg.FileList) == 1 && cfg.FileList[0] == "." {
		cfg.FileList = []string{"/root"}
	}

	var stdout, stderr bytes.Buffer
	a := &app.App{
		Config:   cfg,
		Writer:   &stdout,
		Logger:   log.New(&stderr, "", 0),
		FS:       vfs,
		ExitCode: cli.CodeOk,
	}
	a.Run()

	return runResult{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: a.ExitCode,
	}
}

// ansiRE matches ANSI escape sequences including 24-bit color.
var ansiRE = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// stripANSI removes color codes from s.
func stripANSI(s string) string {
	return ansiRE.ReplaceAllString(s, "")
}

// stripTrailingWS removes trailing whitespace from each line. Internal column
// padding is preserved so alignment regressions still surface.
func stripTrailingWS(s string) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	return strings.Join(lines, "\n")
}

// normalize applies the default normalization for golden comparison.
func normalize(s string) string {
	return stripTrailingWS(stripANSI(s))
}

// goldenDir is the directory under tests/ where golden files live.
const goldenDir = "testdata"

// usedGoldenNames tracks (lowercased) names written this run so that
// case-insensitive filesystems (macOS APFS, Windows NTFS) do not silently
// share a single file between two tests.
var usedGoldenNames = map[string]string{}

// assertGolden compares actual against the golden file at golden/<name>.txt.
// With -update, the golden is rewritten from actual.
func assertGolden(t *testing.T, name, actual string) {
	t.Helper()
	if !strings.HasSuffix(name, ".txt") {
		name += ".txt"
	}
	if prior, ok := usedGoldenNames[strings.ToLower(name)]; ok && prior != name {
		t.Fatalf("golden %q collides with %q on case-insensitive filesystems — use distinct names", name, prior)
	}
	usedGoldenNames[strings.ToLower(name)] = name

	path := filepath.Join(goldenDir, name)
	got := normalize(actual)

	if *updateGolden {
		if err := os.MkdirAll(goldenDir, 0o755); err != nil {
			t.Fatalf("mkdir golden: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v (run `go test ./tests -update` to create)", path, err)
	}
	if normalize(string(want)) != got {
		t.Errorf("golden %s mismatch\n--- want ---\n%s\n--- got ---\n%s",
			name, string(want), got)
	}
}

// assertExitCode fails the test if the run's exit code is not want.
func assertExitCode(t *testing.T, want cli.ExitCode, got cli.ExitCode) {
	t.Helper()
	if got != want {
		t.Errorf("exit code: want %d, got %d", want, got)
	}
}

// assertNoLogs fails if stderr is non-empty.
func assertNoLogs(t *testing.T, r runResult) {
	t.Helper()
	if strings.TrimSpace(r.Stderr) != "" {
		t.Errorf("expected no logger output, got: %q", r.Stderr)
	}
}

// assertContainsLine fails if no line of out matches re.
func assertContainsLine(t *testing.T, out string, re string) {
	t.Helper()
	pat := regexp.MustCompile(re)
	if slices.ContainsFunc(strings.Split(normalize(out), "\n"), pat.MatchString) {
		return
	}
	t.Errorf("expected line matching %q in output:\n%s", re, normalize(out))
}

// assertContains fails if out does not contain substr (after normalization).
func assertContains(t *testing.T, out, substr string) {
	t.Helper()
	if !strings.Contains(normalize(out), substr) {
		t.Errorf("expected %q in output:\n%s", substr, normalize(out))
	}
}

// assertNotContains fails if out contains substr (after normalization).
func assertNotContains(t *testing.T, out, substr string) {
	t.Helper()
	if strings.Contains(normalize(out), substr) {
		t.Errorf("expected %q NOT in output:\n%s", substr, normalize(out))
	}
}

// lines returns non-empty trailing-trimmed lines of normalize(s).
func lines(s string) []string {
	out := []string{}
	for l := range strings.SplitSeq(normalize(s), "\n") {
		if l == "" {
			continue
		}
		out = append(out, l)
	}
	return out
}

// reverseLines returns the lines in reverse order.
func reverseLines(ls []string) []string {
	out := make([]string, len(ls))
	for i, l := range ls {
		out[len(ls)-1-i] = l
	}
	return out
}
