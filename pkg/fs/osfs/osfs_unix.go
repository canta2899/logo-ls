//go:build !windows

package osfs

// All per-file metadata extraction (inode, hardlinks, mode bits, owner/
// group, blocks, xattr) is handled by internal/inspect/platform via a
// single fi.Sys().(*syscall.Stat_t) read. This file exists so the build
// tag pairs with osfs_windows.go.
