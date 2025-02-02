package app

import (
	"os"
	"testing"

	"github.com/canta2899/logo-ls/model"
	"github.com/pborman/getopt/v2"
)

// reset os.Args and the getopt state before calling GetConfigFromCli.
func parseArgs(args []string) *Config {
	// Save original os.Args so we can restore it later.
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	os.Args = args
	getopt.Reset() // Reset the flag parser between tests.
	return GetConfigFromCli()
}

// Verifies that when no file arguments are provided,
// the config defaults to using the current directory.
func TestNoArgs(t *testing.T) {
	cfg := parseArgs([]string{"app"})
	if len(cfg.FileList) != 1 || cfg.FileList[0] != "." {
		t.Errorf("expected FileList to be [\".\"], got %v", cfg.FileList)
	}
	if cfg.SortMode != model.SortAlphabetical {
		t.Errorf("expected default SortMode to be SortAlphabetical, got %v", cfg.SortMode)
	}
	if cfg.LongListingMode != model.LongListingNone {
		t.Errorf("expected default LongListingMode to be LongListingNone, got %v", cfg.LongListingMode)
	}
}

// Verifies that file arguments are added to the config.
func TestFileArgs(t *testing.T) {
	cfg := parseArgs([]string{"app", "file1", "file2"})
	if len(cfg.FileList) != 2 {
		t.Fatalf("expected FileList length 2, got %d", len(cfg.FileList))
	}
	if cfg.FileList[0] != "file1" || cfg.FileList[1] != "file2" {
		t.Errorf("expected FileList to be [\"file1\", \"file2\"], got %v", cfg.FileList)
	}
}

// Verifies that the sort flags set the appropriate sort mode.
func TestSortFlags(t *testing.T) {
	tests := []struct {
		args     []string
		expected model.SortMode
	}{
		{[]string{"app", "-U"}, model.SortNone},
		{[]string{"app", "-v"}, model.SortNatural},
		{[]string{"app", "-X"}, model.SortExtension},
		{[]string{"app", "-t"}, model.SortModTime},
		{[]string{"app", "-S"}, model.SortSize},
	}
	for _, tt := range tests {
		cfg := parseArgs(tt.args)
		if cfg.SortMode != tt.expected {
			t.Errorf("for args %v, expected SortMode %v, got %v", tt.args, tt.expected, cfg.SortMode)
		}
	}
}

// Verifies that the long listing flags are parsed correctly.
func TestLongListingFlags(t *testing.T) {
	tests := []struct {
		args     []string
		expected model.Listing
	}{
		{[]string{"app", "-o"}, model.LongListingOwner},
		{[]string{"app", "-g"}, model.LongListingGroup},
		{[]string{"app", "-l"}, model.LongListingDefault},
	}
	for _, tt := range tests {
		cfg := parseArgs(tt.args)
		if cfg.LongListingMode != tt.expected {
			t.Errorf("for args %v, expected LongListingMode %v, got %v", tt.args, tt.expected, cfg.LongListingMode)
		}
	}
}

// Verifies various boolean flags.
func TestIncludeFlags(t *testing.T) {
	cfg := parseArgs([]string{"app", "-a", "-r", "-R", "-D", "-e", "-i", "-1", "-d", "-G", "-h", "-s", "-T"})
	if cfg.AllMode != model.IncludeAll {
		t.Errorf("expected AllMode to be IncludeAll, got %v", cfg.AllMode)
	}
	if !cfg.Reverse {
		t.Error("expected Reverse to be true")
	}
	if !cfg.Recursive {
		t.Error("expected Recursive to be true")
	}
	if !cfg.GitStatus {
		t.Error("expected GitStatus to be true")
	}
	if !cfg.DisableIcon {
		t.Error("expected DisableIcon to be true")
	}
	if !cfg.ShowInodeNumber {
		t.Error("expected ShowInodeNumber to be true")
	}
	if !cfg.OneFilePerLine {
		t.Error("expected OneFilePerLine to be true")
	}
	if !cfg.Directory {
		t.Error("expected Directory to be true")
	}
	if !cfg.NoGroup {
		t.Error("expected NoGroup to be true")
	}
	if !cfg.HumanReadable {
		t.Error("expected HumanReadable to be true")
	}
	if !cfg.ShowBlockSize {
		t.Error("expected ShowBlockSize to be true")
	}
	if cfg.TimeFormatter == nil {
		t.Error("expected TimeFormatter to be non-nil")
	}
}

// Verifies that the "-A" flag sets the IncludeAlmost mode.
func TestIncludeAlmost(t *testing.T) {
	cfg := parseArgs([]string{"app", "-A"})
	if cfg.AllMode != model.IncludeAlmost {
		t.Errorf("expected AllMode to be IncludeAlmost, got %v", cfg.AllMode)
	}
}
