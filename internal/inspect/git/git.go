// Package git parses `git status --porcelain -z` output and provides a small
// abstraction (Porcelain) for talking to the git binary so the parser can be
// tested in isolation with canned bytes.
package git

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Porcelain is the minimal interface the StatusReader needs to talk to git.
type Porcelain interface {
	// Root returns the absolute path to the top-level of the repository
	// containing dir, or an error if dir is not inside a git repo.
	Root(dir string) (string, error)
	// Status returns the raw bytes of `git status --porcelain -z` for the
	// repository rooted at root.
	Status(root string) ([]byte, error)
}

// StatusReader resolves and caches per-repository status maps.
type StatusReader struct {
	porcelain Porcelain
	// cache keys are repository roots.
	cache map[string]map[string]string
}

// NewStatusReader returns a fresh status reader using the given porcelain
// implementation. The cache is per-instance: there is no global state.
func NewStatusReader(p Porcelain) *StatusReader {
	return &StatusReader{
		porcelain: p,
		cache:     make(map[string]map[string]string),
	}
}

// Status returns a map of absolute-path -> single-char status code for the
// repository containing dir. Returns nil with no error when dir is not inside
// a git repository.
func (r *StatusReader) Status(dir string) map[string]string {
	root, err := r.porcelain.Root(dir)
	if err != nil {
		return nil
	}
	if cached, ok := r.cache[root]; ok {
		return cached
	}
	raw, err := r.porcelain.Status(root)
	if err != nil {
		return nil
	}
	m := ParsePorcelain(root, raw)
	r.cache[root] = m
	return m
}

// StatusRelative returns r.Status(dir) re-keyed to paths relative to dir.
// Returns nil if dir is not in a git repository.
func (r *StatusReader) StatusRelative(dir string) map[string]string {
	repoMap := r.Status(dir)
	if repoMap == nil {
		return nil
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil
	}
	absDir = filepath.Clean(absDir)

	out := make(map[string]string, len(repoMap))
	for absFile, code := range repoMap {
		rel, ok := strings.CutPrefix(absFile, absDir)
		if !ok {
			continue
		}
		rel = strings.TrimPrefix(rel, string(filepath.Separator))
		out[rel] = code
	}
	return out
}

// ParsePorcelain converts the output of `git status --porcelain -z` into a
// map of absolute-path -> status character. Parent directories of every
// changed file are marked with "M" so directory listings can show that they
// contain modifications.
//
// The porcelain v1 format is `XY <space> path`, where XY is exactly two
// characters wide. Records are separated by NUL when -z is used.
func ParsePorcelain(repoRoot string, raw []byte) map[string]string {
	records := splitPorcelain(string(raw))
	result := make(map[string]string, len(records))
	for _, rec := range records {
		if len(rec) < 4 {
			continue
		}
		xy := rec[:2]
		path := rec[3:]
		statusChar := ExtractStatusChar(xy)
		absFile := filepath.Clean(filepath.Join(repoRoot, filepath.FromSlash(path)))
		result[absFile] = statusChar
		for _, parent := range parentDirsWithinRepo(repoRoot, absFile) {
			if !strings.HasSuffix(parent, string(filepath.Separator)) {
				parent += string(filepath.Separator)
			}
			if _, exists := result[parent]; !exists {
				result[parent] = "M"
			}
		}
	}
	return result
}

// splitPorcelain splits a porcelain -z stream into one record per change.
//
// Most records are XY-space-path; renames and copies are XY-space-newpath
// followed by a NUL and the old path, which is consumed but not emitted as
// its own record.
func splitPorcelain(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, "\000")
	var out []string
	for i := 0; i < len(parts); i++ {
		p := parts[i]
		if p == "" {
			continue
		}
		if len(p) >= 2 && (p[0] == 'R' || p[0] == 'C' || p[1] == 'R' || p[1] == 'C') {
			// Skip the followed "old path" record.
			i++
		}
		out = append(out, p)
	}
	return out
}

// ExtractStatusChar picks a single status character from a porcelain XY code.
// "??" (untracked) is mapped to "U" and "!!" (ignored) is mapped to "I" to
// match the rendering conventions used elsewhere in logo-ls.
func ExtractStatusChar(xy string) string {
	xy = strings.TrimSpace(xy)
	if xy == "" {
		return "?"
	}
	if xy == "??" {
		return "U"
	}
	if xy == "!!" {
		return "I"
	}
	for _, r := range xy {
		if r != ' ' && r != '\t' {
			return string(r)
		}
	}
	return "?"
}

func parentDirsWithinRepo(repoRoot, absFile string) []string {
	var parents []string
	dir := filepath.Clean(filepath.Dir(absFile))
	for dir != repoRoot {
		if dir == "" || dir == string(filepath.Separator) {
			break
		}
		parents = append(parents, dir)
		next := filepath.Dir(dir)
		if next == dir {
			break
		}
		dir = next
	}
	return parents
}

// ExecPorcelain shells out to the real `git` binary.
type ExecPorcelain struct{}

func (ExecPorcelain) Root(dir string) (string, error) {
	out, err := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (ExecPorcelain) Status(root string) ([]byte, error) {
	return exec.Command("git", "-C", root, "status", "--porcelain", "--ignored", "-z").Output()
}
