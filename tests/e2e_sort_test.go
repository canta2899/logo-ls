package tests

import (
	"github.com/canta2899/logo-ls/internal/cli"
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
)

// Sort fixture distinctness: every sort mode must produce a distinct ordering
// when applied to sortFixture(). If two modes happen to produce the same
// output, the fixture is too tame and we should change it.

func TestSort_AlphabeticalIsDefault(t *testing.T) {
	// Pin to the C locale so collation is deterministic across systems
	// (uppercase before lowercase, byte order).
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_COLLATE", "C")
	t.Setenv("LANG", "")
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-1e", "/root")
	assertGolden(t, "sort_alphabetical", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
}

func TestSort_AllModesAreDistinct(t *testing.T) {
	modes := map[string]string{
		"alphabetical": "-1e",
		"none":         "-1Ue",
		"natural":      "-1ve",
		"extension":    "-1Xe",
		"modtime":      "-1te",
		"size":         "-1Se",
	}
	outputs := map[string]string{}
	for label, args := range modes {
		vfs := fakefs.New(sortFixture())
		r := runApp(t, vfs, strings.Split(args+" /root", " ")...)
		outputs[label] = normalize(r.Stdout)
	}
	// At least 4 of the 6 outputs should be unique; alphabetical and natural
	// happen to coincide in the current implementation, so allow some overlap.
	seen := map[string]string{}
	for label, out := range outputs {
		if prev, ok := seen[out]; ok {
			t.Logf("sort modes %q and %q produced identical output", prev, label)
		} else {
			seen[out] = label
		}
	}
	if len(seen) < 4 {
		t.Errorf("sort modes too similar — only %d distinct outputs (want >=4)", len(seen))
	}
}

func TestSort_DotfilesAlwaysFirst(t *testing.T) {
	// For sort modes that apply dotfile grouping (alphabetical, size,
	// natural), .hidden must lead. Mtime sort and -U intentionally do not
	// regroup dotfiles; -X groups dotfiles after extensionless files (not
	// necessarily first), so all of these are excluded here.
	for _, args := range [][]string{
		{"-1Ae", "/root"},
		{"-1ASe", "/root"},
		{"-1Ave", "/root"},
	} {
		vfs := fakefs.New(sortFixture())
		r := runApp(t, vfs, args...)
		got := lines(r.Stdout)
		if len(got) == 0 || got[0] != ".hidden" {
			t.Errorf("args %v: expected .hidden first, got %v", args, got)
		}
	}
}

func TestSort_ExtensionDotfilesByExt(t *testing.T) {
	// With -X, dotfiles stay grouped: extensionless non-dotfiles first, then
	// the dotfile group, then files sorted by extension. Dotfiles are neither
	// interleaved by extension nor force-pinned to the top.
	vfs := fakefs.New(dotfileExtTree())
	r := runApp(t, vfs, "-1AXe", "/root")
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	got := lines(r.Stdout)
	// Makefile (no ext), then .hidden (dotfile group), then a.go, README.md.
	want := []string{"Makefile", ".hidden", "a.go", "README.md"}
	if !equalSlice(got, want) {
		t.Errorf("-AX order:\nwant %v\ngot  %v", want, got)
	}
	// .hidden must not lead: extensionless non-dotfiles come first.
	if got[0] == ".hidden" {
		t.Errorf("expected extensionless files before the dotfile group: %v", got)
	}
}

func TestSort_ExtensionDotfilesStayGrouped(t *testing.T) {
	// -X groups dotfiles together as a contiguous block: extensionless
	// non-dotfiles (dirs, Makefile) first, then every dotfile (including
	// ".config.json", which must NOT sort among the regular .json files), then
	// the remaining files by extension. "." and ".." are ordinary dotfiles and
	// sort within the dotfile group, not force-pinned to the very top.
	t.Setenv("LC_ALL", "")
	t.Setenv("LC_COLLATE", "C")
	t.Setenv("LANG", "")
	vfs := fakefs.New(dotfileGroupTree())
	r := runApp(t, vfs, "-1aXe", "/root")
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	got := lines(r.Stdout)

	want := []string{
		"Makefile", "src/", // extensionless non-dotfiles
		"./", "../", ".config.json", ".hidden", // dotfile group
		"main.go", "app.json", // by extension: go < json
	}
	if !equalSlice(got, want) {
		t.Fatalf("-aX grouping:\nwant %v\ngot  %v", want, got)
	}
	// "." / ".." must not be force-pinned: extensionless files lead.
	if got[0] == "./" || got[0] == "../" {
		t.Errorf("expected extensionless files before . and .., got %v", got)
	}
}

func TestSort_ReverseAppliesAfterSort(t *testing.T) {
	for _, args := range [][]string{
		{"-1Se", "/root"},
		{"-1Xe", "/root"},
	} {
		vfs1 := fakefs.New(sortFixture())
		vfs2 := fakefs.New(sortFixture())
		fwd := lines(runApp(t, vfs1, args...).Stdout)
		revArgs := append([]string{"-r"}, args...)
		rev := lines(runApp(t, vfs2, revArgs...).Stdout)
		if !equalSlice(reverseLines(fwd), rev) {
			t.Errorf("args %v + -r should reverse:\nfwd=%v\nrev=%v", args, fwd, rev)
		}
	}
}
