package inspect

import "github.com/canta2899/logo-ls/internal/icons"

type defaultIconResolver struct{}

func (defaultIconResolver) Resolve(name, ext, indicator string) *icons.IconInfo {
	return icons.Resolve(name, ext, indicator)
}

// DefaultIconResolver returns the package-level icon resolver.
func DefaultIconResolver() IconResolver { return defaultIconResolver{} }
