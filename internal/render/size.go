package render

import "fmt"

// formatSize renders b as either a raw byte count or a human-readable
// 1K/2.3M-style size when humanReadable is true.
func formatSize(b int64, humanReadable bool) string {
	if !humanReadable {
		return fmt.Sprintf("%d", b)
	}
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	size := float64(b) / float64(div)
	if size == float64(int64(size)) {
		return fmt.Sprintf("%d%c", int64(size), "KMGTPE"[exp])
	}
	return fmt.Sprintf("%.1f%c", size, "KMGTPE"[exp])
}
