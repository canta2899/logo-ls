package format

import (
	"github.com/canta2899/logo-ls/internal/inspect/git"
)

// defaultReader is a process-wide StatusReader backed by the real git binary.
// It is retained as a small convenience for code that does not yet plumb a
// per-instance reader through (osfs.GitStatus). Subsequent phases of the
// refactor migrate this responsibility onto the inspector.
var defaultReader = git.NewStatusReader(git.ExecPorcelain{})

// GetFilesGitStatus returns a map of paths-relative-to-p -> status code, or
// nil if p is not inside a git repository.
func GetFilesGitStatus(p string) map[string]string {
	return defaultReader.StatusRelative(p)
}
