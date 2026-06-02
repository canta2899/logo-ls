//go:build windows

package osfs

// All per-file metadata extraction (inode, hardlinks, mode bits, owner/
// group, blocks, xattr) is handled by internal/inspect/platform via a
// platform-specific reader. This file exists so the build tag pairs with
// osfs_unix.go.
