// this file contain dir type definition

//go:build windows
// +build windows

package model

import "os"

func DirBlocks(info *file, fi os.FileInfo) {
}

func GetOwnerGroupInfo(fi os.FileInfo, noGroup bool, longListingMode model.Listing) (o string, g string) {
	return
}
