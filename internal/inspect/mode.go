package inspect

import (
	"fmt"
	iofs "io/fs"
	"strconv"
	"strings"
)

// ModeString builds the long-mode 11-char permission string from raw mode
// bits, sticky bits and an HasXAttr flag. Returns a placeholder if mode is
// zero (i.e. the inspector failed to stat).
func ModeString(mode iofs.FileMode, sticky, stickyX, hasXAttr bool) string {
	if mode == 0 {
		return fmt.Sprintf("%-*s", 11, strings.Repeat("?", 11))
	}
	s := mode.String()
	// FileMode.String() returns "L..." for symlinks; the conventional ls
	// rendering is lowercase 'l'.
	if mode&iofs.ModeSymlink != 0 && len(s) > 0 && s[0] == 'L' {
		s = "l" + s[1:]
	}
	if sticky {
		// FileMode.String() encodes ModeSticky as 't' at position 9 for
		// Unix-style modes; substitute the trailing char to honour stickyX
		// vs sticky-only.
		if len(s) >= 10 {
			if stickyX {
				s = s[:9] + "t" + s[10:]
			} else {
				s = s[:9] + "T" + s[10:]
			}
		}
	}
	if hasXAttr {
		s += "@"
	}
	return fmt.Sprintf("%-*s", 11, s)
}

// FormatHardLinks renders a hard-link count for the long listing.
func FormatHardLinks(n uint64) string { return strconv.FormatUint(n, 10) }
