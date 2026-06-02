package inspect

import "github.com/canta2899/logo-ls/internal/icons"

type defaultIconResolver struct {
	ext *icons.Extension
}

func (r defaultIconResolver) Resolve(name, ext, indicator string) *icons.IconInfo {
	return icons.ResolveWith(r.ext, name, ext, indicator)
}

// DefaultIconResolver returns the package-level icon resolver with no overrides.
func DefaultIconResolver() IconResolver { return defaultIconResolver{} }

// IconResolverWith returns a resolver that consults the given user extension
// before falling back to the built-in icon tables.
func IconResolverWith(ext *icons.Extension) IconResolver {
	return defaultIconResolver{ext: ext}
}
