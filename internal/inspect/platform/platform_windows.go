//go:build windows

package platform

import (
	"strconv"

	"github.com/canta2899/logo-ls/fs"
	"golang.org/x/sys/windows"
)

// SysProvider is the optional interface a FileInfo can implement to surface
// platform data without going through OS-specific stat. fakefs uses it.
type SysProvider interface {
	PlatformStat() Stat
}

type winReader struct{}

func newPlatformReader() Reader { return &winReader{} }

func (winReader) Read(absPath string, fi fs.FileInfo, opts Options) Stat {
	if fi == nil {
		return Stat{}
	}
	if sp, ok := fi.(SysProvider); ok {
		return sp.PlatformStat()
	}
	if absPath == "" {
		return Stat{}
	}
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
		return Stat{}
	}
	defer windows.CloseHandle(handle)

	var info windows.ByHandleFileInformation
	if err := windows.GetFileInformationByHandle(handle, &info); err != nil {
		return Stat{}
	}
	inode := (uint64(info.FileIndexHigh) << 32) | uint64(info.FileIndexLow)
	return Stat{
		Inode:     strconv.FormatUint(inode, 10),
		HardLinks: uint64(info.NumberOfLinks),
	}
}

func (winReader) LookupOwner(uid uint32) string { return "" }
func (winReader) LookupGroup(gid uint32) string { return "" }
