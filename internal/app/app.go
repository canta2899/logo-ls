// Package app defines the main application logic
package app

import (
	"bytes"
	"fmt"
	"io"
	iofs "io/fs"
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/canta2899/logo-ls/internal/cli"
	"github.com/canta2899/logo-ls/internal/icons"
	"github.com/canta2899/logo-ls/internal/inspect"
	"github.com/canta2899/logo-ls/internal/inspect/git"
	"github.com/canta2899/logo-ls/internal/render"
	isort "github.com/canta2899/logo-ls/internal/sort"
	"github.com/canta2899/logo-ls/pkg/fs"
)

const cannotAccessFmt = "cannot access %q: %v\n"

type App struct {
	Config   *cli.Config
	Writer   io.Writer
	ExitCode cli.ExitCode
	Logger   *log.Logger
	FS       fs.FS
	// GitReader is optional; when nil the app falls back to FS.GitStatus
	GitReader    *git.StatusReader
	IconOverride *icons.Override
}

// gitStatusFor returns the status map for dir, using the per-app reader when
// configured and falling back to the legacy FS.GitStatus otherwise.
func (a *App) gitStatusFor(dir string) map[string]string {
	if a.GitReader != nil {
		return a.GitReader.StatusRelative(dir)
	}
	return a.FS.GitStatus(dir)
}

type Args struct {
	Files []FileEntry
	Dirs  []DirectoryEntry
}

type RecursiveLookupFrame struct {
	entry  *DirectoryEntry
	header string // if non-empty, printed as: "\n<icon><header>:\n"
}

func (a *App) Exit() {
	os.Exit(int(a.ExitCode))
}

// Write panics if copying buf to the writer fails.
func (a *App) Write(buf *bytes.Buffer) {
	if _, err := io.Copy(a.Writer, buf); err != nil {
		panic(err)
	}
}

func (a *App) GetArguments() *Args {
	slices.Sort(a.Config.FileList)

	args := &Args{}

	for _, argPath := range a.Config.FileList {
		abs, err := a.FS.Abs(argPath)
		if err != nil {
			a.Logger.Printf("cannot get absolute path for %q: %v\n", argPath, err)
			a.ExitCode.SetSerious()
			continue
		}

		f, err := a.FS.Open(abs)
		if err != nil {
			a.Logger.Printf(cannotAccessFmt, argPath, err)
			a.ExitCode.SetSerious()
			continue
		}

		fi, err := f.Stat()
		if err != nil {
			a.Logger.Printf("cannot stat %q: %v\n", argPath, err)
			f.Close()
			a.ExitCode.SetSerious()
			continue
		}

		if fi.IsDir() {
			args.Dirs = append(args.Dirs, DirectoryEntry{
				File:    f,
				AbsPath: abs,
			})
		} else {
			args.Files = append(args.Files, FileEntry{
				Info:    fi,
				AbsPath: abs,
			})
			f.Close() // no need to keep file handles for single files
		}
	}
	return args
}

func (a *App) Run() {
	args := a.GetArguments()

	if len(args.Files) > 0 {
		filesDir := a.ProcessFiles(args.Files)
		a.PrintDirectory(filesDir)
		if len(args.Dirs) > 0 {
			fmt.Fprintln(a.Writer)
		}
	}

	if a.Config.Recursive {
		a.processDirsRecursively(args.Dirs)
	} else {
		a.processDirsNonRecursively(args.Dirs)
	}
}

func (a *App) processDirsRecursively(dirs []DirectoryEntry) {
	currentAbs, _ := a.FS.Abs(".")
	openDirIcon := OpenDirIconString(!a.Config.DisableIcon)

	for i, dirEntry := range dirs {
		if i > 0 {
			fmt.Fprintln(a.Writer)
		}

		relName := dirEntry.Name()
		if rel, err := a.FS.Rel(currentAbs, dirEntry.Name()); err == nil {
			relName = rel
		}

		fmt.Fprintf(a.Writer, "%s:\n", openDirIcon+relName)

		a.recurseDirectory(&dirEntry, currentAbs)
	}
}

func (a *App) processDirsNonRecursively(dirs []DirectoryEntry) {
	pName := len(dirs) > 1
	openDirIcon := OpenDirIconString(!a.Config.DisableIcon)

	for i, dirEntry := range dirs {
		if pName {
			fmt.Fprintf(a.Writer, "%s:\n", openDirIcon+dirEntry.Name())
		}

		d, err := a.ProcessDirectory(&dirEntry)
		dirEntry.Close()
		if err != nil {
			a.Logger.Printf(cannotAccessFmt, dirEntry.Name(), err)
			a.ExitCode.SetSerious()
		}

		a.PrintDirectory(d)
		if i < len(dirs)-1 {
			fmt.Fprintln(a.Writer)
		}
	}
}

func (a *App) recurseDirectory(start *DirectoryEntry, startingAbsolutePath string) {
	stack := []*RecursiveLookupFrame{{entry: start, header: ""}}

	for len(stack) > 0 {
		idx := len(stack) - 1
		current := stack[idx]
		stack = stack[:idx]

		if current.header != "" {
			fmt.Fprintf(a.Writer, "\n%s:\n", OpenDirIconString(!a.Config.DisableIcon)+current.header)
		}

		d, err := a.ProcessDirectory(current.entry)
		current.entry.Close()
		if err != nil {
			a.Logger.Printf(cannotAccessFmt, current.entry.Name(), err)
			a.ExitCode.SetMinor()
		}

		a.PrintDirectory(d)

		if d == nil || len(d.Dirs) == 0 {
			continue
		}

		sort.Strings(d.Dirs)
		stack = a.pushSubdirFrames(stack, current.entry, d.Dirs, startingAbsolutePath)
	}
}

func (a *App) pushSubdirFrames(stack []*RecursiveLookupFrame, parent *DirectoryEntry, dirs []string, startingAbsolutePath string) []*RecursiveLookupFrame {
	for i := len(dirs) - 1; i >= 0; i-- {
		frame := a.openSubdirFrame(parent, dirs[i], startingAbsolutePath)
		if frame == nil {
			continue
		}
		stack = append(stack, frame)
	}
	return stack
}

func (a *App) openSubdirFrame(parent *DirectoryEntry, subdir, startingAbsolutePath string) *RecursiveLookupFrame {
	subdirFullPath := a.FS.Join(parent.Name(), subdir)

	childPath := subdirFullPath
	if rel, err := a.FS.Rel(startingAbsolutePath, subdirFullPath); err == nil {
		childPath = rel
	}

	f, err := a.FS.Open(subdirFullPath)
	if err != nil {
		a.Logger.Printf(cannotAccessFmt, childPath, err)
		a.ExitCode.SetMinor()
		return nil
	}
	abs, err := a.FS.Abs(subdirFullPath)
	if err != nil {
		a.Logger.Println("Cannot compute abs path for:", childPath)
		f.Close()
		return nil
	}
	return &RecursiveLookupFrame{
		entry:  &DirectoryEntry{File: f, AbsPath: abs},
		header: childPath,
	}
}

func (a *App) ProcessFiles(files []FileEntry) *Directory {
	t := new(Directory)
	isLong := a.Config.LongListingMode != cli.LongListingNone

	for _, fileEntry := range files {
		entry := a.buildEntry(fileEntry.AbsPath, fileEntry.Info, isLong)
		t.Files = append(t.Files, entry)
	}

	return t
}

func (a *App) ProcessDirectory(d *DirectoryEntry) (*Directory, error) {
	defer func() {
		_ = d.Close()
	}()

	dirStat, err := d.File.Stat()
	if err != nil {
		return nil, err
	}

	dirModel, err := a.populateDirectory(d, dirStat)
	return dirModel, err
}

func (a *App) populateDirectory(d *DirectoryEntry, dirStat fs.FileInfo) (*Directory, error) {
	t := new(Directory)
	isLong := a.Config.LongListingMode != cli.LongListingNone

	a.maybeAttachSelfEntry(t, d, dirStat, isLong)

	if a.Config.Directory {
		t.Files = append(t.Files, t.Info)
		return t, nil
	}

	entries, err := d.File.ReadDir(0)
	// proceed even on error: entries may contain a partial list

	var gitRepoStatus map[string]string
	if a.Config.GitStatus {
		gitRepoStatus = a.gitStatusFor(d.Name())
	}

	showHidden := a.Config.AllMode != cli.IncludeDefault

	for _, de := range entries {
		name := de.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		a.appendChildEntry(t, d, de, gitRepoStatus, isLong)
	}

	if a.Config.AllMode == cli.IncludeAll {
		a.appendDotEntries(t, d, isLong)
	}
	return t, err
}

func (a *App) maybeAttachSelfEntry(t *Directory, d *DirectoryEntry, dirStat fs.FileInfo, isLong bool) {
	if a.Config.AllMode != cli.IncludeAll && !a.Config.Directory {
		return
	}
	t.Info = a.buildEntry(d.Name(), dirStat, isLong)
	if !a.Config.Directory {
		t.Info.Name = "."
		t.Info.Base = "."
		t.Info.Ext = ""
	}
	if !a.Config.DisableIcon {
		t.Info.Icon = icons.OpenDir()
	}
}

func (a *App) appendChildEntry(t *Directory, d *DirectoryEntry, de fs.DirEntry, gitRepoStatus map[string]string, isLong bool) {
	name := de.Name()
	fullpath := a.FS.Join(d.Name(), name)
	fi, infoErr := de.Info() // fi might be nil on error
	if infoErr != nil {
		a.Logger.Printf(cannotAccessFmt, fullpath, infoErr)
		a.ExitCode.SetMinor()
	}

	entry := a.buildEntry(fullpath, fi, isLong)

	if fi == nil {
		a.fillMissingInfo(entry, de)
	}

	if gitRepoStatus != nil {
		entry.GitStatus = gitRepoStatus[name+a.FS.Separator()]
		if entry.GitStatus == "" {
			entry.GitStatus = gitRepoStatus[name]
		}
	}

	t.Files = append(t.Files, entry)
	if de.IsDir() {
		t.Dirs = append(t.Dirs, name+"/")
	}
}

func (a *App) fillMissingInfo(entry *inspect.InspectedEntry, de fs.DirEntry) {
	if de.IsDir() {
		entry.Indicator = "/"
	} else if de.Type()&iofs.ModeSymlink != 0 {
		entry.Indicator = "@"
	}
	if !a.Config.DisableIcon {
		entry.Icon = icons.ResolveWith(a.IconOverride, entry.Base, entry.Ext, entry.Indicator)
	}
}

func (a *App) appendDotEntries(t *Directory, d *DirectoryEntry, isLong bool) {
	if t.Info != nil {
		t.Files = append(t.Files, t.Info)
	}

	pp := a.FS.Dir(d.Name())
	pStat, _ := a.FS.Lstat(pp)

	parentEntry := a.buildEntry(pp, pStat, isLong)
	parentEntry.Name = ".."
	parentEntry.Base = ".."
	parentEntry.Ext = ""

	if !a.Config.DisableIcon {
		parentEntry.Icon = icons.OpenDir()
	}

	t.Files = append(t.Files, parentEntry)
	t.Parent = parentEntry
}

// inspectorFor builds an Inspector for exactly the columns the current mode needs.
func (a *App) inspectorFor(isLong bool) *inspect.Inspector {
	showOwner := a.Config.LongListingMode == cli.LongListingDefault ||
		a.Config.LongListingMode == cli.LongListingOwner
	showGroup := !a.Config.NoGroup &&
		(a.Config.LongListingMode == cli.LongListingDefault ||
			a.Config.LongListingMode == cli.LongListingGroup)
	return inspect.New(a.FS, inspect.IconResolverWith(a.IconOverride), inspect.Options{
		Long:            isLong,
		ShowOwner:       showOwner,
		ShowGroup:       showGroup,
		ShowInode:       a.Config.ShowInodeNumber,
		ShowBlocks:      a.Config.ShowBlockSize,
		ResolveSymlinks: !a.Config.DisableIcon,
		DisableIcon:     a.Config.DisableIcon,
	})
}

// buildEntry inspects fullPath. When fi is nil, returns a stub entry with just the name set.
func (a *App) buildEntry(fullPath string, fi fs.FileInfo, isLong bool) *inspect.InspectedEntry {
	insp := a.inspectorFor(isLong)
	return insp.Inspect(fullPath, fi)
}

func (a *App) PrintDirectory(d *Directory) {
	if d == nil {
		return
	}
	isort.Sort(d.Files, a.Config.SortMode, a.Config.Reverse)
	render.Render(a.Writer, d.Files, render.Options{
		Mode:          a.renderMode(),
		ShowIcon:      !a.Config.DisableIcon,
		ShowInode:     a.Config.ShowInodeNumber,
		ShowBlocks:    a.Config.ShowBlockSize,
		HumanReadable: a.Config.HumanReadable,
		TimeFormatter: a.Config.TimeFormatter,
	})
}

func (a *App) renderMode() render.Mode {
	switch {
	case a.Config.LongListingMode != cli.LongListingNone:
		return render.ModeLong
	case a.Config.OneFilePerLine:
		return render.ModeOneFilePerLine
	default:
		return render.ModeShort
	}
}
