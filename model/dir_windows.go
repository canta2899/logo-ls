//go:build windows
// +build windows

package model

import (
	"os"
)

func DirBlocks(info *Entry, fi os.FileInfo) {
}

func GetOwnerGroupInfo(fi os.FileInfo, noGroup bool, longListingMode Listing) (string, string) {
	return "", ""
}
