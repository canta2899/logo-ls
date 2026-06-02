package inspect

import (
	"fmt"
	iofs "io/fs"
	"strconv"
	"strings"

	"github.com/canta2899/logo-ls/model"
)

// ToLegacy converts an InspectedEntry into a *model.Entry so the existing
// renderer keeps working unchanged during the migration. Once the renderer
// consumes InspectedEntry directly (Phase 8), this can be deleted.
func ToLegacy(e *InspectedEntry, modeStr string, owner, group string, p Pather) *model.Entry {
	out := &model.Entry{}
	name, ext := splitNameExt(e.Name, p)
	out.Name = name
	out.Ext = ext
	out.Size = e.Size
	out.ModTime = e.ModTime
	out.Indicator = e.Indicator
	out.Mode = modeStr
	out.ModeBits = uint32(e.Mode)
	out.NumHardLinks = e.HardLinks
	out.Owner = owner
	out.Group = group
	out.Blocks = e.Blocks
	out.Icon = e.Icon
	out.InodeNumber = e.Inode
	out.GitStatus = e.GitStatus
	return out
}

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
		// FileMode.String() already encodes ModeSticky as 't' at position 9
		// for Unix-style modes, but it uses uppercase 'T' if other-exec is
		// not set. Replace position 9 to honour StickyX explicitly.
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

// FormatHardLinks renders a count for legacy rows.
func FormatHardLinks(n uint64) string { return strconv.FormatUint(n, 10) }
