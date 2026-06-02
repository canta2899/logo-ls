//go:build !windows

package platform

import (
	"os/user"
	"strconv"
	"syscall"

	"github.com/canta2899/logo-ls/pkg/fs"
	"golang.org/x/sys/unix"
)

// SysProvider is the optional interface a FileInfo can implement to surface
// platform data without going through *syscall.Stat_t. fakefs uses this so
// tests don't have to fabricate raw stat structs.
type SysProvider interface {
	PlatformStat() Stat
}

type unixReader struct {
	users  map[uint32]string
	groups map[uint32]string
}

func newPlatformReader() Reader {
	return &unixReader{
		users:  make(map[uint32]string),
		groups: make(map[uint32]string),
	}
}

func (r *unixReader) Read(absPath string, fi fs.FileInfo, opts Options) Stat {
	if fi == nil {
		return Stat{}
	}
	// fakefs supplies a SysProvider instead of a real *syscall.Stat_t.
	if sp, ok := fi.(SysProvider); ok {
		s := sp.PlatformStat()
		if opts.WantXAttr {
			s.HasXAttr = listXAttr(absPath)
		}
		return s
	}
	st, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return Stat{}
	}
	s := Stat{
		Inode:     strconv.FormatUint(uint64(st.Ino), 10),
		HardLinks: uint64(st.Nlink),
		Blocks:    int64(st.Blocks),
		UID:       st.Uid,
		GID:       st.Gid,
		Sticky:    st.Mode&unix.S_ISVTX != 0,
	}
	if s.Sticky && st.Mode&unix.S_IXOTH != 0 {
		s.StickyX = true
	}
	if opts.WantXAttr {
		s.HasXAttr = listXAttr(absPath)
	}
	return s
}

func listXAttr(path string) bool {
	count, err := unix.Listxattr(path, nil)
	return err == nil && count > 0
}

func (r *unixReader) LookupOwner(uid uint32) string {
	if n, ok := r.users[uid]; ok {
		return n
	}
	if u, err := user.LookupId(strconv.FormatUint(uint64(uid), 10)); err == nil {
		r.users[uid] = u.Username
		return u.Username
	}
	r.users[uid] = ""
	return ""
}

func (r *unixReader) LookupGroup(gid uint32) string {
	if n, ok := r.groups[gid]; ok {
		return n
	}
	if g, err := user.LookupGroupId(strconv.FormatUint(uint64(gid), 10)); err == nil {
		r.groups[gid] = g.Name
		return g.Name
	}
	r.groups[gid] = ""
	return ""
}
