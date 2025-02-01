//go:build windows
// +build windows

package format

import (
	"os"
	"strconv"

	"golang.org/x/sys/windows"
)

// GetInodeNumber returns the Windows "File ID" (FileIndexHigh/Low) as a string.
// This is the closest approximation to a Unix inode number on NTFS and other FS.
func GetInodeNumber(path string) string {
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
	err = windows.GetFileInformationByHandle(handle, &fileInfo)
	if err != nil {
		return ""
	}

	// Combine FileIndexHigh and FileIndexLow into one 64-bit value:
	inode := (uint64(fileInfo.FileIndexHigh) << 32) | uint64(fileInfo.FileIndexLow)

	return strconv.FormatUint(inode, 10)
}

// GetHardLinkCount returns the number of hard links that a file or directory on
// Windows %has, by calling GetFileInformationByHandle on the file handle.
//
// Note that it returns 0 on any error (e.g., if the file cannot be opened).
func GetHardLinkCount(absPath string) uint64 {
	// CreateFileW (from "golang.org/x/sys/windows") opens the file and returns a handle.
	handle, err := windows.CreateFile(
		windows.StringToUTF16Ptr(absPath),
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

	// Retrieve the file (or directory) information.
	var fileInfo windows.ByHandleFileInformation
	err = windows.GetFileInformationByHandle(handle, &fileInfo)
	if err != nil {
		return 0
	}

	// The NumberOfLinks field is a WORD (16-bit) in the BY_HANDLE_FILE_INFORMATION struct.
	return uint64(fileInfo.NumberOfLinks)
}

func GetModeExtended(fi *os.FileInfo, fullPath string) string {
	return (*fi).Mode().String()
}
