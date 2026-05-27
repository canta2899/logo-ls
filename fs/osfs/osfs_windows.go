//go:build windows

package osfs

import (
	"strconv"

	"github.com/canta2899/logo-ls/fs"
	"golang.org/x/sys/windows"
)

func (o *osFS) InodeNumber(path string) string {
	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(path),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(handle)

	var fileInfo windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &fileInfo); err != nil {
		return ""
	}
	inode := (uint64(fileInfo.FileIndexHigh) << 32) | uint64(fileInfo.FileIndexLow)
	return strconv.FormatUint(inode, 10)
}

func (o *osFS) HardLinks(path string) uint64 {
	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(path),
		windows.GENERIC_READ,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	if err != nil {
		return 0
	}
	defer windows.CloseHandle(handle)

	var fileInfo windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &fileInfo); err != nil {
		return 0
	}
	return uint64(fileInfo.NumberOfLinks)
}

func (o *osFS) ModeExtended(fi fs.FileInfo, path string) string {
	if fi == nil {
		return "???????????"
	}
	return fi.Mode().String()
}

func (o *osFS) OwnerGroup(fi fs.FileInfo, showOwner, showGroup bool) (string, string) {
	return "", ""
}

func (o *osFS) Blocks(fi fs.FileInfo) int64 {
	return 0
}
