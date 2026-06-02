package cli

import "time"

// Timestamp renders a *time.Time for the long-mode timestamp column.
// Different implementations select the locale/precision.
type Timestamp interface {
	Format(t *time.Time) string
}

// DefaultFormatter mirrors `ls -l`: "Jan 02 15:04" for entries from the
// current year, "Jan 02  2006" otherwise.
type DefaultFormatter struct{}

func (*DefaultFormatter) Format(t *time.Time) string {
	if t.Year() == time.Now().Year() {
		return t.Format("Jan 02 15:04")
	}
	return t.Format("Jan 02  2006")
}

// ExtendedFormatter renders the full date and time, used by `-T`.
type ExtendedFormatter struct{}

func (*ExtendedFormatter) Format(t *time.Time) string {
	return t.Format("Jan 02 15:04:05 2006")
}

// GetFormatter returns the formatter selected by the -T flag.
func GetFormatter(extended bool) Timestamp {
	if extended {
		return &ExtendedFormatter{}
	}
	return &DefaultFormatter{}
}
