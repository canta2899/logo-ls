//go:build !windows

package osfs

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/canta2899/logo-ls/fs"
	"golang.org/x/sys/unix"
)

var (
	uidCache = make(map[string]string)
	gidCache = make(map[string]string)
)

func (o *osFS) InodeNumber(path string) string {
	fi, err := os.Stat(path)
	if err != nil {
		return ""
	}
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return ""
	}
	return strconv.Itoa(int(stat.Ino))
}

func (o *osFS) HardLinks(path string) uint64 {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return 0
	}
	return uint64(stat.Nlink)
}

func (o *osFS) ModeExtended(fi fs.FileInfo, path string) string {
	if fi == nil {
		return fmt.Sprintf("%-*s", 11, "???????????")
	}
	mode := fi.Mode()
	modeStr := mode.String()

	var stat unix.Stat_t
	if err := unix.Stat(path, &stat); err == nil {
		if stat.Mode&unix.S_ISVTX != 0 {
			if stat.Mode&unix.S_IXOTH != 0 {
				modeStr = modeStr[:9] + "t"
			} else {
				modeStr = modeStr[:9] + "T"
			}
		}
	}

	count, err := unix.Listxattr(path, nil)
	if err == nil && count > 0 {
		modeStr += "@"
	}

	return fmt.Sprintf("%-*s", 11, modeStr)
}

func (o *osFS) OwnerGroup(fi fs.FileInfo, showOwner, showGroup bool) (string, string) {
	if fi == nil {
		return "", ""
	}
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return "", ""
	}
	var owner, group string
	if showOwner {
		uid := strconv.Itoa(int(stat.Uid))
		if n, ok := uidCache[uid]; ok {
			owner = n
		} else if u, err := user.LookupId(uid); err == nil {
			owner = u.Username
			uidCache[uid] = owner
		}
	}
	if showGroup {
		gid := strconv.Itoa(int(stat.Gid))
		if n, ok := gidCache[gid]; ok {
			group = fmt.Sprintf(" %v  ", n)
		} else if g, err := user.LookupGroupId(gid); err == nil {
			group = fmt.Sprintf(" %v  ", g.Name)
			gidCache[gid] = g.Name
		}
	}
	return owner, group
}

func (o *osFS) Blocks(fi fs.FileInfo) int64 {
	if fi == nil {
		return 0
	}
	if s, ok := fi.Sys().(*syscall.Stat_t); ok {
		return s.Blocks
	}
	return 0
}
