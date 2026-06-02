// Package platform isolates OS-specific extraction of per-file metadata.
// On Unix it reads directly from fi.Sys().(*syscall.Stat_t) so a single Lstat
// produces the inode, link count, block count, ownership and sticky bits.
// On Windows it returns best-effort zero values; the renderer omits those
// columns where the OS doesn't provide them.
package platform

import (
	"github.com/canta2899/logo-ls/fs"
)

// Stat carries every per-file scalar the inspector needs from a single
// underlying stat. Strings are kept as raw values; padding is the renderer's
// concern.
type Stat struct {
	Inode     string
	HardLinks uint64
	Blocks    int64
	UID       uint32
	GID       uint32
	Sticky    bool // S_ISVTX
	StickyX   bool // sticky AND other-executable (controls `t` vs `T`)
	HasXAttr  bool // populated only when WantXAttr is set and the OS supports it
}

// Options controls expensive optional lookups.
type Options struct {
	// WantXAttr asks the platform layer to call Listxattr (Unix only).
	WantXAttr bool
}

// Reader is the platform-specific reader. On Unix it uses a sentinel that
// optionally accepts an absolute path so it can call Listxattr; on systems
// without xattr support the path is ignored.
type Reader interface {
	// Read extracts metadata from fi. absPath is only used for xattr lookup
	// when opts.WantXAttr is true.
	Read(absPath string, fi fs.FileInfo, opts Options) Stat
	// LookupOwner returns the username for uid, cached for the lifetime of
	// the reader. Returns "" if not resolvable.
	LookupOwner(uid uint32) string
	// LookupGroup returns the group name for gid, cached.
	LookupGroup(gid uint32) string
}

// NewReader returns a platform-appropriate reader with fresh per-instance
// caches. There is no package-global state.
func NewReader() Reader { return newPlatformReader() }
