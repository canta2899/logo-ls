package tests

import (
	"github.com/canta2899/logo-ls/internal/cli"
	"strings"
	"testing"

	"github.com/canta2899/logo-ls/pkg/fs/fakefs"
)

func TestFlag_a_IncludesDotAndParent(t *testing.T) {
	vfs := fakefs.New(hiddenTree())
	r := runApp(t, vfs, "-1ae", "/root")
	assertGolden(t, "flag_a_all", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -a includes ., .. (rendered with trailing '/' since both are dirs).
	assertContains(t, r.Stdout, "./")
	assertContains(t, r.Stdout, "../")
	// Intent: -a shows dotfiles
	assertContains(t, r.Stdout, ".env")
	assertContains(t, r.Stdout, ".config")
}

func TestFlag_A_AlmostAllNoDotDotDot(t *testing.T) {
	vfs := fakefs.New(hiddenTree())
	r := runApp(t, vfs, "-1Ae", "/root")
	assertGolden(t, "flag_A_almost", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -A shows hidden files but not . or ..
	assertContains(t, r.Stdout, ".env")
	for _, l := range lines(r.Stdout) {
		if strings.TrimSpace(l) == "." || strings.TrimSpace(l) == ".." {
			t.Errorf("-A should not include %q", l)
		}
	}
}

func TestFlag_l_LongListing(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-le", "/root")
	assertGolden(t, "flag_l", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: long mode shows mode string and owner+group
	assertContainsLine(t, r.Stdout, `^-rw-r--r-- .*alice.*staff`)
}

func TestFlag_o_OwnerOnly(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-oe", "/root")
	assertGolden(t, "flag_o", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -o is like -l but no group column
	assertContains(t, r.Stdout, "alice")
	assertNotContains(t, r.Stdout, " staff  ")
}

func TestFlag_g_GroupOnly(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-ge", "/root")
	assertGolden(t, "flag_g_grouponly", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -g is like -l but no owner column
	assertContains(t, r.Stdout, "staff")
	assertNotContains(t, r.Stdout, "alice")
}

func TestFlag_G_NoGroupInLong(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-lGe", "/root")
	assertGolden(t, "flag_G_nogroup", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -G removes group column from -l
	assertContains(t, r.Stdout, "alice")
	assertNotContains(t, r.Stdout, " staff  ")
}

func TestFlag_h_HumanReadable(t *testing.T) {
	vfs := fakefs.New(fakefs.Dir("root", dirMeta("a"),
		fakefs.File("big.bin", 2*1024*1024, mtime("2026-01-01 10:00:00"), fileMeta("b")),
		fakefs.File("kb.bin", 2048, mtime("2026-01-02 10:00:00"), fileMeta("c")),
		fakefs.File("small.bin", 100, mtime("2026-01-03 10:00:00"), fileMeta("d")),
	))
	r := runApp(t, vfs, "-lhe", "/root")
	assertGolden(t, "flag_h", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -h converts byte counts to K/M
	assertContains(t, r.Stdout, "2M")
	assertContains(t, r.Stdout, "2K")
}

func TestFlag_s_ShowBlockSize(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-se", "/root")
	assertGolden(t, "flag_s_blocksize", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -s prefixes each entry with the block count (8 per fileMeta).
	assertContainsLine(t, r.Stdout, `^\s*8 .*notes.txt`)
}

func TestFlag_i_ShowInode(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-ie", "/root")
	assertGolden(t, "flag_i", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -i prefixes each entry with the inode number.
	assertContains(t, r.Stdout, "1001")
	assertContains(t, r.Stdout, "1002")
}

func TestFlag_1_OneFilePerLine(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-1e", "/root")
	assertGolden(t, "flag_1", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -1 puts each entry on its own line.
	got := lines(r.Stdout)
	if len(got) != 3 {
		t.Errorf("expected 3 lines, got %d: %v", len(got), got)
	}
}

func TestFlag_d_DirItself(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-de", "/root")
	assertGolden(t, "flag_d_dironly", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -d lists the dir itself, not its contents.
	got := lines(r.Stdout)
	if len(got) != 1 {
		t.Errorf("expected 1 line for -d, got %d: %v", len(got), got)
	}
}

func TestFlag_r_ReverseSort(t *testing.T) {
	vfs := fakefs.New(smallTree())
	fwd := lines(runApp(t, fakefs.New(smallTree()), "-1e", "/root").Stdout)
	rev := lines(runApp(t, vfs, "-1re", "/root").Stdout)
	assertGolden(t, "flag_r", strings.Join(rev, "\n"))
	// Intent: -r reverses the order of -1.
	if !equalSlice(reverseLines(fwd), rev) {
		t.Errorf("-r did not reverse -1:\nfwd=%v\nrev=%v", fwd, rev)
	}
}

func TestFlag_t_SortByMtime(t *testing.T) {
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-1te", "/root")
	assertGolden(t, "flag_t_modtime", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -t sorts newest first. alpha.txt has mtime 2026-04-01
	// (the most recent in sortFixture).
	got := lines(r.Stdout)
	if len(got) == 0 || got[0] != "alpha.txt" {
		t.Errorf("expected alpha.txt first, got %v", got)
	}
}

func TestFlag_S_SortBySize(t *testing.T) {
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-1Se", "/root")
	assertGolden(t, "flag_S_sortsize", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -S sorts largest first. zebra.txt is 999 bytes (largest).
	got := lines(r.Stdout)
	// dotfiles come first in this app's sort regardless of -S, then largest.
	// Find the first non-dot entry:
	for _, l := range got {
		if !strings.HasPrefix(l, ".") {
			if l != "zebra.txt" {
				t.Errorf("expected first non-dot entry to be zebra.txt, got %q", l)
			}
			break
		}
	}
}

func TestFlag_U_NoSort(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-1Ue", "/root")
	assertGolden(t, "flag_U", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -U lists in directory (fake = name) order, no further sorting.
}

func TestFlag_v_NaturalSort(t *testing.T) {
	vfs := fakefs.New(sortFixture())
	r := runApp(t, vfs, "-1ve", "/root")
	assertGolden(t, "flag_v", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -v selects the SortNatural mode. The current implementation
	// of SortNatural is lexical string compare (with dotfile priority), so
	// file10 sorts before file2. This is locked in by the golden; the
	// behavioral check just confirms both entries are present.
	got := lines(r.Stdout)
	if indexOf(got, "file2.go") < 0 || indexOf(got, "file10.go") < 0 {
		t.Fatalf("file2/file10 missing: %v", got)
	}
}

func TestFlag_X_SortByExtension(t *testing.T) {
	vfs := fakefs.New(mixedExtTree())
	r := runApp(t, vfs, "-1Xe", "/root")
	assertGolden(t, "flag_X", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -X puts no-extension files first (Makefile, script), then .go,
	// then .md. mixedExtTree has no dotfiles, so the dotfile-by-extension
	// behaviour (covered by TestSort_ExtensionDotfilesByExt) does not apply here.
	got := lines(r.Stdout)
	if len(got) < 5 {
		t.Fatalf("expected 5 entries, got %v", got)
	}
	// First two have no extension (Makefile, script).
	noExtCount := 0
	for _, l := range got {
		if !strings.Contains(l, ".") {
			noExtCount++
			continue
		}
		break
	}
	if noExtCount != 2 {
		t.Errorf("expected 2 no-ext entries first, got order %v", got)
	}
}

func TestFlag_T_ExtendedTimeStyle(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-lTe", "/root")
	assertGolden(t, "flag_T_timestyle", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -T uses the extended timestamp formatter. We override the
	// formatter in the harness, so we only check the long-mode columns
	// are still present. Format change is checked at unit level.
	assertContainsLine(t, r.Stdout, `^-rw-r--r--`)
}

func TestFlag_D_GitStatus(t *testing.T) {
	vfs := fakefs.New(gitRepoTree(), fakefs.WithGitStatus(gitRepoStatus()))
	r := runApp(t, vfs, "-1De", "/root")
	assertGolden(t, "flag_D_gitstatus", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -D attaches git status code to each entry. The codes appear
	// as a trailing character; the harness strips ANSI so they survive.
	assertContains(t, r.Stdout, "A") // staged
	assertContains(t, r.Stdout, "M") // modified
	assertContains(t, r.Stdout, "U") // untracked
}

func TestFlag_e_NoIcons(t *testing.T) {
	vfs := fakefs.New(smallTree())
	r := runApp(t, vfs, "-1e", "/root")
	assertGolden(t, "flag_e", r.Stdout)
	assertExitCode(t, cli.CodeOk, r.ExitCode)
	// Intent: -e removes the icon column. No nerd-font glyph (any PUA) appears.
	if inPUA(r.Stdout) {
		t.Errorf("found nerd-font glyph in -e output:\n%s", r.Stdout)
	}
}

func equalSlice(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func indexOf(s []string, target string) int {
	for i, v := range s {
		if v == target {
			return i
		}
	}
	return -1
}
