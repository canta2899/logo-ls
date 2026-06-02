package inspect

import (
	"github.com/canta2899/logo-ls/format"
	"github.com/canta2899/logo-ls/icons"
)

// defaultIconResolver delegates to format.GetIcon. Once the icon resolver
// is fully extracted in a later phase, this shim goes away.
type defaultIconResolver struct{}

func (defaultIconResolver) Resolve(name, ext, indicator string) *icons.IconInfo {
	return format.GetIcon(name, ext, indicator)
}

// DefaultIconResolver returns an IconResolver backed by the current
// format.GetIcon rules.
func DefaultIconResolver() IconResolver { return defaultIconResolver{} }
