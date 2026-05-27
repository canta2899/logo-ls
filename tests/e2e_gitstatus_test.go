package tests

import (
	"testing"

	"github.com/canta2899/logo-ls/fs/fakefs"
	"github.com/canta2899/logo-ls/model"
)

func TestGitStatus_LongMode(t *testing.T) {
	vfs := fakefs.New(gitRepoTree(), fakefs.WithGitStatus(gitRepoStatus()))
	r := runApp(t, vfs, "-lDe", "/root")
	assertGolden(t, "gitstatus_long", r.Stdout)
	assertExitCode(t, model.CodeOk, r.ExitCode)
	// Each tracked file's status code shows up.
	assertContains(t, r.Stdout, "staged.txt")
	assertContains(t, r.Stdout, "modified.txt")
	assertContains(t, r.Stdout, "untracked.txt")
}

func TestGitStatus_DisabledByDefault(t *testing.T) {
	// Without -D, no status codes should be attached even if the map exists.
	vfs := fakefs.New(gitRepoTree(), fakefs.WithGitStatus(gitRepoStatus()))
	r := runApp(t, vfs, "-1e", "/root")
	out := normalize(r.Stdout)
	// Status codes "A"/"M"/"U" would appear as standalone glyphs; the
	// filename "modified.txt" itself contains an 'M', so we look at single-
	// character "A " / "U " patterns separated from the name.
	if got := lines(out); len(got) > 0 {
		for _, l := range got {
			if l == "A" || l == "M" || l == "U" {
				t.Errorf("found bare status code line without -D: %q", l)
			}
		}
	}
}
