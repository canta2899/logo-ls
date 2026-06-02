package inspect

import (
	iofs "io/fs"
	"os"
	"strings"

	"github.com/canta2899/logo-ls/fs"
	"github.com/canta2899/logo-ls/icons"
	"github.com/canta2899/logo-ls/internal/inspect/platform"
)

// foldHome rewrites paths inside HOME to use "~" so symlink targets in long
// mode read cleanly (matches the legacy osfs behaviour).
func foldHome(p string) string {
	home := os.Getenv("HOME")
	if home == "" {
		return p
	}
	if p == home {
		return "~"
	}
	if strings.HasPrefix(p, home+string(os.PathSeparator)) {
		return "~" + strings.TrimPrefix(p, home)
	}
	return p
}

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

// namedOwner is implemented by FileInfo values that know their owner/group
// without needing a uid/gid lookup. fakefs uses it so tests don't depend on
// the host's user database.
type namedOwner interface {
	OwnerName() string
	GroupName() string
}

// Inspector is the single place that touches fs.FS for per-file metadata.
// Owns per-instance caches that used to live as package globals.
type Inspector struct {
	fs       fs.FS
	icons    IconResolver
	options  Options
	platform platform.Reader
}

// New returns a fresh inspector.
func New(filesystem fs.FS, ir IconResolver, opts Options) *Inspector {
	return &Inspector{
		fs:       filesystem,
		icons:    ir,
		options:  opts,
		platform: platform.NewReader(),
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

	wantStat := i.options.ShowInode || i.options.ShowBlocks || i.options.Long
	if wantStat {
		stat := i.platform.Read(absPath, fi, platform.Options{WantXAttr: i.options.Long && i.options.WantXAttr})
		if i.options.ShowInode {
			e.Inode = stat.Inode
		}
		if i.options.ShowBlocks {
			e.Blocks = stat.Blocks
			if e.Blocks == 0 && e.Kind == KindFile {
				e.Blocks = (e.Size + 511) / 512
			}
		}
		if i.options.Long {
			e.HardLinks = stat.HardLinks
			if e.HardLinks == 0 {
				e.HardLinks = 1
			}
			e.HasXAttr = stat.HasXAttr
			e.Sticky = stat.Sticky
			e.StickyX = stat.StickyX
			if no, ok := fi.(namedOwner); ok {
				if i.options.ShowOwner {
					e.Owner = no.OwnerName()
				}
				if i.options.ShowGroup {
					e.Group = no.GroupName()
				}
			} else {
				if i.options.ShowOwner {
					e.Owner = i.platform.LookupOwner(stat.UID)
				}
				if i.options.ShowGroup {
					e.Group = i.platform.LookupGroup(stat.GID)
				}
			}
		}
	}

	if e.Kind == KindSymlink && (i.options.ResolveSymlinks || i.options.Long) {
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

	e.Indicator = i.indicatorFor(e, fi.Mode())

	if !i.options.DisableIcon && i.icons != nil {
		name, ext := splitNameExt(e.Name, i.fs)
		if e.Kind == KindSymlink && e.LinkResolved != nil {
			tname, text := splitNameExt(e.LinkResolved.Name, i.fs)
			tind := classifyIndicator(e.LinkResolved.Mode)
			e.Icon = i.icons.Resolve(tname, text, tind)
		} else {
			e.Icon = i.icons.Resolve(name, ext, e.Indicator)
		}
	}

	return e
}

// indicatorFor returns the trailing classifier glyph ("/", "@", "*", ...).
// For symlinks in long mode it emits " ~> target" using the inspected
// LinkTarget; the HOME->~ folding is done here so it doesn't rely on a
// per-FS Indicator method. A symlink whose target can't be resolved gets
// an empty indicator in long mode (matches the previous behaviour).
func (i *Inspector) indicatorFor(e *InspectedEntry, m iofs.FileMode) string {
	if m&iofs.ModeSymlink != 0 {
		if i.options.Long {
			if e.LinkTarget == "" {
				return ""
			}
			return " ~> " + foldHome(e.LinkTarget)
		}
		return "@"
	}
	return classifyIndicator(m)
}

// classifyIndicator returns the indicator string for non-symlinks (and the
// short-mode classification of any entry).
func classifyIndicator(m iofs.FileMode) string {
	switch {
	case m&iofs.ModeDir != 0:
		return "/"
	case m&iofs.ModeNamedPipe != 0:
		return "|"
	case m&iofs.ModeSymlink != 0:
		return "@"
	case m&iofs.ModeSocket != 0:
		return "="
	case m&0o111 != 0:
		return "*"
	}
	return ""
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
