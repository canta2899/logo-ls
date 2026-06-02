package inspect

import (
	"strings"

	"github.com/canta2899/logo-ls/fs"
	"github.com/canta2899/logo-ls/icons"
)

// Options configures what the inspector collects per entry. The renderer
// declares its needs and the inspector skips work it doesn't have to do.
type Options struct {
	Long            bool // populate Mode/Owner/Group/HardLinks
	ShowOwner       bool
	ShowGroup       bool
	ShowInode       bool
	ShowBlocks      bool
	WantXAttr       bool // call Listxattr; only meaningful in long mode
	ResolveSymlinks bool // populate LinkResolved (long mode follows targets for icons)
	DisableIcon     bool
}

// IconResolver picks an icon for an entry. Defaults to the package-level
// resolver that uses the existing format.GetIcon rules.
type IconResolver interface {
	Resolve(name, ext, indicator string) *icons.IconInfo
}

// GitStatusReader returns a git status map keyed by absolute path, or nil
// if the directory is not inside a repository.
type GitStatusReader interface {
	Status(dir string) map[string]string
}

// Inspector is the single place that touches fs.FS for per-file metadata.
// Owns per-instance caches that used to live as package globals.
type Inspector struct {
	fs      fs.FS
	icons   IconResolver
	options Options
}

// New returns a fresh inspector.
func New(filesystem fs.FS, ir IconResolver, opts Options) *Inspector {
	return &Inspector{
		fs:      filesystem,
		icons:   ir,
		options: opts,
	}
}

// Inspect builds an InspectedEntry for absPath. The caller passes the
// FileInfo it already has from ReadDir/Lstat/Stat so we don't repeat the
// syscall here.
func (i *Inspector) Inspect(absPath string, fi fs.FileInfo) *InspectedEntry {
	e := &InspectedEntry{AbsPath: absPath}
	if fi == nil {
		e.Name = i.fs.Base(absPath)
		return e
	}

	e.Name = fi.Name()
	e.Mode = fi.Mode()
	e.Size = fi.Size()
	e.ModTime = fi.ModTime()
	e.Kind = kindFromMode(fi.Mode())

	if i.options.ShowInode {
		e.Inode = i.fs.InodeNumber(absPath)
	}
	if i.options.ShowBlocks {
		e.Blocks = i.fs.Blocks(fi)
	}

	if i.options.Long {
		e.HardLinks = i.fs.HardLinks(absPath)
		owner, group := i.fs.OwnerGroup(fi, i.options.ShowOwner, i.options.ShowGroup)
		e.Owner = owner
		e.Group = group
	}

	e.Indicator = i.fs.Indicator(absPath, i.options.Long)

	if e.Kind == KindSymlink && i.options.ResolveSymlinks {
		if target, err := i.fs.EvalSymlinks(absPath); err == nil {
			e.LinkTarget = target
			if tfi, terr := i.fs.Stat(target); terr == nil {
				e.LinkResolved = &InspectedEntry{
					AbsPath: target,
					Name:    tfi.Name(),
					Mode:    tfi.Mode(),
					Size:    tfi.Size(),
					ModTime: tfi.ModTime(),
					Kind:    kindFromMode(tfi.Mode()),
				}
			}
		}
	}

	if !i.options.DisableIcon && i.icons != nil {
		name, ext := splitNameExt(e.Name, i.fs)
		if e.Kind == KindSymlink && e.LinkResolved != nil {
			tname, text := splitNameExt(e.LinkResolved.Name, i.fs)
			tind := i.fs.Indicator(e.LinkTarget, i.options.Long)
			e.Icon = i.icons.Resolve(tname, text, tind)
		} else {
			e.Icon = i.icons.Resolve(name, ext, e.Indicator)
		}
	}

	return e
}

// SplitNameExt splits a filename into its (name, ext) parts, treating
// dotfiles without an interior dot as having no extension. Exported so the
// renderer can use the same rule when rendering legacy adapter output.
func SplitNameExt(name string, p Pather) (string, string) {
	return splitNameExt(name, p)
}

// Pather is the subset of fs.FS used for path manipulation in this package.
type Pather interface {
	Ext(path string) string
}

func splitNameExt(name string, p Pather) (string, string) {
	if strings.HasPrefix(name, ".") && !strings.Contains(name[1:], ".") {
		return name, ""
	}
	ext := p.Ext(name)
	return name[0 : len(name)-len(ext)], ext
}
