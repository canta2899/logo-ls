package format

import "time"

type Timestamp interface {
	Format(t *time.Time) string
}

type DefaultFormatter struct{}
type ExtendedFormatter struct{}

func (*DefaultFormatter) Format(t *time.Time) string {
	if t.Year() == time.Now().Year() {
		return t.Format("Jan 02 15:04")
	}

	return t.Format("Jan 02  2006")
}

func (*ExtendedFormatter) Format(t *time.Time) string {
	return t.Format("Jan 02 15:04:05 2006")
}

func GetFormatter(extended bool) Timestamp {
	if extended {
		return &ExtendedFormatter{}
	}

	return &DefaultFormatter{}
}
