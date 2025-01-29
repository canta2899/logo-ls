//go:build !windows
// +build !windows

package format

import (
	"os"
	"strconv"
	"syscall"
)

func GetInodeNumber(path string) string {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return ""
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)

	if !ok {
		return ""
	}

	return strconv.Itoa(int(stat.Ino))
}

func GetHardLinkCount(absPath string) uint64 {
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return 0
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return 0
	}

	return uint64(stat.Nlink)
}
