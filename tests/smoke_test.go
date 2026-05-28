package tests

import (
	"testing"

	"github.com/canta2899/logo-ls/fs/fakefs"
	"github.com/canta2899/logo-ls/model"
)

func TestHarnessSmoke(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-1e", "/root")
	assertExitCode(t, model.CodeOk, r.ExitCode)
	assertNoLogs(t, r)
	if r.Stdout == "" {
		t.Fatal("empty stdout")
	}
}
