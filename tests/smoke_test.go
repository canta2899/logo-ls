package tests

import (
	"github.com/canta2899/logo-ls/internal/cli"
	"testing"

	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
)

func TestHarnessSmoke(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-1e", "/root")
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	assertNoLogs(t, r)
	if r.Stdout == "" {
		t.Fatal("empty stdout")
	}
}
