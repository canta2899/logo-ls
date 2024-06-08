// this file contain dir type definition

//go:build !windows
// +build !windows

package model

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

var grpMap = make(map[string]string)
var userMap = make(map[string]string)

func DirBlocks(info *Entry, fi os.FileInfo) {
	if s, ok := fi.Sys().(*syscall.Stat_t); ok {
		info.Blocks = s.Blocks
	}
}

func GetOwnerGroupInfo(fi os.FileInfo, noGroup bool, longListingMode Listing) (o string, g string) {

	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		if longListingMode == LongListingDefault || longListingMode == LongListingOwner {
			UID := strconv.Itoa(int(stat.Uid))
			if n, ok := userMap[UID]; ok {
				o = n
			} else {
				u, err := user.LookupId(UID)
				if err != nil {
					o = ""
				} else {
					o = u.Username
					userMap[UID] = u.Username
				}
			}
		}

		if !noGroup && (longListingMode == LongListingDefault || longListingMode == LongListingGroup) {
			GID := strconv.Itoa(int(stat.Gid))
			if n, ok := grpMap[GID]; ok {
				g = n
			} else {
				grp, err := user.LookupGroupId(GID)
				if err != nil {
					g = ""
				} else {
					g = grp.Name
					grpMap[GID] = grp.Name
				}
			}
		}
	}

	return
}

func getFileBlocks(fi os.FileInfo) int64 {
	if s, ok := fi.Sys().(*syscall.Stat_t); ok {
		return s.Blocks
	}
	return 0
}
