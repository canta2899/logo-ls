//go:build windows
// +build windows

package format

func GetInodeNumber(path string) string {
	return ""
}

func GetHardLinkCount(absPath string) uint64 {
	return 0
}
