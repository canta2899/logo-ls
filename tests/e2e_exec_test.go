package tests

import (
	"testing"

	"github.com/canta2899/logo-ls/fs/fakefs"
	"github.com/canta2899/logo-ls/model"
)

// Exec-bit tests: executable files render with the '*' indicator in -F-like
// modes (which is the default of logo-ls).

func TestExec_OneLineStar(t *testing.T) {
	vfs := fakefs.New(execTree())
	r := runApp(t, vfs, "-1e", "/root")
	assertGolden(t, "exec_oneline", r.Stdout)
	assertExitCode(t, model.CodeOk, r.ExitCode)
	assertContains(t, r.Stdout, "run.sh*")
	// Non-executable does NOT pick up the star.
	assertNotContains(t, r.Stdout, "regular.txt*")
}
