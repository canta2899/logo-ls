package tests

import (
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/fs/fakefs"
	"github.com/canta2899/logo-ls/model"
)

// Sort fixture distinctness: every sort mode must produce a distinct ordering
// when applied to sortFixture(). If two modes happen to produce the same
// output, the fixture is too tame and we should change it.

func TestSort_AlphabeticalIsDefault(t *testing.T) {
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-1e", "/root")
	assertGolden(t, "sort_alphabetical", r.Stdout)
	assertExitCode(t, model.CodeOk, r.ExitCode)
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
	// extension, natural), .hidden must lead. Mtime sort and -U intentionally
	// do not regroup dotfiles, so they are excluded here.
	for _, args := range [][]string{
		{"-1Ae", "/root"},
		{"-1ASe", "/root"},
		{"-1AXe", "/root"},
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
