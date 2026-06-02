package inspect

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/canta2899/logo-ls/model"
)

// ToLegacy converts an InspectedEntry into a *model.Entry so the existing
// renderer keeps working unchanged during the migration. Once the renderer
// consumes InspectedEntry directly (Phase 5+), this can be deleted.
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

// LegacyModeString builds the long-mode permission string for legacy
// renderers. Falls back to a placeholder when raw mode is zero.
func LegacyModeString(e *InspectedEntry) string {
	if e == nil {
		return strings.Repeat("?", 11)
	}
	s := e.Mode.String()
	if e.HasXAttr {
		s += "@"
	}
	return fmt.Sprintf("%-*s", 11, s)
}

// FormatHardLinks renders a count for legacy rows.
func FormatHardLinks(n uint64) string { return strconv.FormatUint(n, 10) }
