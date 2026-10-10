package tests

import (
	"bytes"
	"log"
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/internal/app"
	"github.com/canta2899/logo-ls/internal/cli"
	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
)

// Non-tty output (as set up in main) must be plain: no ANSI codes and no
// trailing padding. runApp normalizes both away, so check the raw bytes.
func TestPlainOutput_NoANSINoTrailingSpaces(t *testing.T) {
	for _, flag := range []string{"-1", "-l"} {
		cfg, _, err := cli.BuildConfig([]string{"logo-ls", "--no-override", flag, "/root"})
		if err != nil {
			t.Fatal(err)
		}
		cfg.TimeFormatter = fixedTime{}
		cfg.NoColor = true

		var stdout bytes.Buffer
		a := &app.App{
			Config:   cfg,
			Writer:   &stdout,
			Logger:   log.New(&bytes.Buffer{}, "", 0),
			FS:       fakefs.New(smallTree()),
			ExitCode: cli.CodeOk,
		}
		a.Run()

		out := stdout.String()
		if strings.Contains(out, "\x1b") {
			t.Errorf("%s: unexpected ANSI codes: %q", flag, out)
		}
		for _, l := range strings.Split(out, "\n") {
			if l != strings.TrimRight(l, " ") {
				t.Errorf("%s: trailing spaces in %q", flag, l)
			}
		}
	}
}
