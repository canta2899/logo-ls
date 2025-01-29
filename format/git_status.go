package format

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

var (
	// will store previous git status results to avoid re-running multiple times
	statusCache = make(map[string]map[string]string)

	// will provide syncronization for the statusCache to ensure a single write
	cacheMu sync.Mutex
)

// ComputeGitRepo ensures we have the git status map cached for the repository containing startPath.
// Returns a map of [absolute path -> status code] for all files in the repo, including directories,
// together with the repo root path and an error if not in a git repo or if 'git' fails.
func ComputeGitRepo(startPath string) (map[string]string, string, error) {
	root, err := getGitRoot(startPath)
	if err != nil {
		return nil, "", err // not a Git repo, or error
	}

	cacheMu.Lock()
	defer cacheMu.Unlock()

	if cached, ok := statusCache[root]; ok {
		return cached, root, nil
	}

	repoMap, err := computeStatusMap(root)
	if err != nil {
		return nil, "", err
	}
	statusCache[root] = repoMap
	return repoMap, root, nil
}

// Returns a map of "relative path -> single-letter code" specifically under directory p
// adding the directory markers for each subdirectory.
func GetFilesGitStatus(p string) map[string]string {
	repoMap, _, err := ComputeGitRepo(p)
	if err != nil {
		return nil // not a git repo
	}

	absP, err := filepath.Abs(p)
	if err != nil {
		return nil
	}
	absP = filepath.Clean(absP)

	results := make(map[string]string)
	for absFile, code := range repoMap {
		if strings.HasPrefix(absFile, absP) {
			// produce a relative path
			rel := strings.TrimPrefix(absFile, absP)
			rel = strings.TrimPrefix(rel, string(filepath.Separator))
			results[rel] = code
		}
	}
	return results
}

// Discards the entire git status cache.
func ClearCache() {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	statusCache = make(map[string]map[string]string)
}

// Runs `git status --porcelain -z` and constructs a map of [absoluteFilePath, code]
// for changed files. Additionally, it **also** populates parent directories with `"●"`.
func computeStatusMap(repoRoot string) (map[string]string, error) {
	out, err := runGitStatusPorcelain(repoRoot)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(out), "\000")
	result := make(map[string]string, len(lines))

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// lines typically look like "?? file" or "M  file" or "R100 new\000old"
		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}
		xy := strings.TrimSpace(parts[0]) // e.g. "??", "M", "A", "R", etc.
		pathPart := parts[1]

		statusChar := extractStatusChar(xy)

		// Convert to an absolute path
		absFilePath := filepath.Join(repoRoot, filepath.FromSlash(pathPart))
		absFilePath = filepath.Clean(absFilePath)

		// Store the actual file status
		result[absFilePath] = statusChar

		// Mark all parent dirs within the repo, up to (but not including) repoRoot
		for _, parentDir := range parentDirsWithinRepo(repoRoot, absFilePath) {
			if !strings.HasSuffix(parentDir, string(filepath.Separator)) {
				parentDir += string(filepath.Separator)
			}
			result[parentDir] = "●"
		}
	}
	return result, nil
}

// Returns a slice of all parent directories of absFilePath that lie withing the root of the repo.
// E.g. if absFilePath = /root/app/file.go, returns ["/root/app", "/root"]. Stops if above root.
func parentDirsWithinRepo(repoRoot, absFilePath string) []string {
	var parents []string

	// Start from the directory portion
	dir := filepath.Dir(absFilePath)
	dir = filepath.Clean(dir)

	// Keep going up until we reach or pass the repoRoot
	for {
		if dir == repoRoot {
			break
		}
		if dir == "" || dir == string(filepath.Separator) {
			break
		}

		parents = append(parents, dir)

		newDir := filepath.Dir(dir)

		if newDir == dir { // Can't go further
			break
		}
		dir = newDir
	}
	return parents
}

// Finds the top-level .git directory via `git -C path rev-parse --show-toplevel`.
func getGitRoot(path string) (string, error) {
	cmd := exec.Command("git", "-C", path, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository: %v", err)
	}
	root := strings.TrimSpace(string(out))
	return root, nil
}

// Calls `git -C root status --porcelain -z` and returns the output as a byte slice.
func runGitStatusPorcelain(root string) ([]byte, error) {
	cmd := exec.Command("git", "-C", root, "status", "--porcelain", "-z")
	return cmd.Output()
}

// Picks a single status character from e.g. "??", "M", " A", etc.
func extractStatusChar(xy string) string {
	xy = strings.TrimSpace(xy)
	if xy == "" {
		return "?" // fallback if we get an empty code
	}

	// If Git porcelain code is "??", interpret as "U"
	if xy == "??" {
		return "U"
	}

	// Return the first non-whitespace character (e.g. 'M', 'A', 'D')
	for _, r := range xy {
		if r != ' ' && r != '\t' {
			return string(r)
		}
	}
	return "?"
}
