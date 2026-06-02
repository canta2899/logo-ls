package inspect

import "github.com/canta2899/logo-ls/internal/icons"

type defaultIconResolver struct {
	ov *icons.Override
}

func (r defaultIconResolver) Resolve(name, ext, indicator string) *icons.IconInfo {
	return icons.ResolveWith(r.ov, name, ext, indicator)
}

// DefaultIconResolver returns the package-level icon resolver with no overrides.
func DefaultIconResolver() IconResolver { return defaultIconResolver{} }

// IconResolverWith returns a resolver that applies the given user overrides
// on top of the built-in icon tables. A nil override yields default behavior.
func IconResolverWith(ov *icons.Override) IconResolver {
	return defaultIconResolver{ov: ov}
}
