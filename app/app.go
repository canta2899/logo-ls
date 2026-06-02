// Package app defines the main application logic
package app

import (
	"bytes"
	"fmt"
	iofs "io/fs"
	"io"
	"log"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/canta2899/logo-ls/format"
	"github.com/canta2899/logo-ls/fs"
	"github.com/canta2899/logo-ls/internal/inspect"
	"github.com/canta2899/logo-ls/internal/inspect/git"
	"github.com/canta2899/logo-ls/internal/render"
	"github.com/canta2899/logo-ls/model"
)

// App represents the main application that holds configuration, a writer, exit codes, and a logger.
type App struct {
	Config   *Config
	Writer   io.Writer
	ExitCode model.ExitCode
	Logger   *log.Logger
	FS       fs.FS
	// GitReader is optional; when nil the app falls back to FS.GitStatus
	// (legacy path) so existing test harnesses keep working.
	GitReader *git.StatusReader
}

// gitStatusFor returns the status map for dir, using the per-app reader when
// configured and falling back to the legacy FS.GitStatus otherwise.
func (a *App) gitStatusFor(dir string) map[string]string {
	if a.GitReader != nil {
		return a.GitReader.StatusRelative(dir)
	}
	return a.FS.GitStatus(dir)
}

// Args stores the parsed command-line arguments as separate files and directories.
type Args struct {
	Files []model.FileEntry
	Dirs  []model.DirectoryEntry
}

type RecursiveLookupFrame struct {
	entry  *model.DirectoryEntry
	header string // if non-empty, printed as: "\n<icon><header>:\n"
}

// Exit terminates the program with the current stored exit code.
func (a *App) Exit() {
	os.Exit(int(a.ExitCode))
}

// Write copies a buffer to the App's writer, panicking on error.
func (a *App) Write(buf *bytes.Buffer) {
	if _, err := io.Copy(a.Writer, buf); err != nil {
		panic(err)
	}
}

// GetArguments parses the configured file list and categorizes them as files or directories.
// It also sets an error exit code for entries that are not accessible.
func (a *App) GetArguments() *Args {
	// Sort all user inputs.
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
			a.Logger.Printf("cannot access %q: %v\n", argPath, err)
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
			args.Dirs = append(args.Dirs, model.DirectoryEntry{
				File:    f,
				AbsPath: abs,
			})
		} else {
			args.Files = append(args.Files, model.FileEntry{
				Info:    fi,
				AbsPath: abs,
			})
			f.Close() // no need to keep file handles for single files
		}
	}
	return args
}

// Run is the main entry point that orchestrates listing files/directories and printing results.
func (a *App) Run() {
	args := a.GetArguments()

	// Process and display all files first.
	if len(args.Files) > 0 {
		filesDir := a.ProcessFiles(args.Files)
		a.PrintDirectory(filesDir)
		if len(args.Dirs) > 0 {
			fmt.Fprintln(a.Writer)
		}
	}

	// Process and display all directories (recursively or not).
	if a.Config.Recursive {
		a.processDirsRecursively(args.Dirs)
	} else {
		a.processDirsNonRecursively(args.Dirs)
	}
}

// Prints each directory and its subdirectories.
func (a *App) processDirsRecursively(dirs []model.DirectoryEntry) {
	currentAbs, _ := a.FS.Abs(".")
	openDirIcon := model.OpenDirIconString(!a.Config.DisableIcon)

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

// Prints each directory (but not subdirectories).
func (a *App) processDirsNonRecursively(dirs []model.DirectoryEntry) {
	pName := len(dirs) > 1
	openDirIcon := model.OpenDirIconString(!a.Config.DisableIcon)

	for i, dirEntry := range dirs {
		if pName {
			fmt.Fprintf(a.Writer, "%s:\n", openDirIcon+dirEntry.Name())
		}

		d, err := a.ProcessDirectory(&dirEntry)
		dirEntry.Close()
		if err != nil {
			a.Logger.Printf("cannot access %q: %v\n", dirEntry.Name(), err)
			a.ExitCode.SetSerious()
		}

		a.PrintDirectory(d)
		if i < len(dirs)-1 {
			fmt.Fprintln(a.Writer)
		}
	}
}

// Processes a directory, prints it, and recurses through subdirectories if any.
func (a *App) recurseDirectory(start *model.DirectoryEntry, startingAbsolutePath string) {
	stack := []*RecursiveLookupFrame{{entry: start, header: ""}}

	for len(stack) > 0 {
		idx := len(stack) - 1
		current := stack[idx]
		stack = stack[:idx]

		if current.header != "" {
			fmt.Fprintf(a.Writer, "\n%s:\n", model.OpenDirIconString(!a.Config.DisableIcon)+current.header)
		}

		d, err := a.ProcessDirectory(current.entry)
		current.entry.Close()
		if err != nil {
			a.Logger.Printf("cannot access %q: %v\n", current.entry.Name(), err)
			a.ExitCode.SetMinor()
		}

		a.PrintDirectory(d)

		if d == nil || len(d.Dirs) == 0 {
			continue
		}

		sort.Strings(d.Dirs)

		for i := len(d.Dirs) - 1; i >= 0; i-- {
			subdir := d.Dirs[i]
			childPath := a.FS.Join(current.entry.Name(), subdir)
			if rel, err := a.FS.Rel(startingAbsolutePath, childPath); err == nil {
				childPath = rel
			}

			subdirFullPath := a.FS.Join(current.entry.Name(), subdir)
			f, err := a.FS.Open(subdirFullPath)
			if err != nil {
				a.Logger.Printf("cannot access %q: %v\n", childPath, err)
				a.ExitCode.SetMinor()
				continue
			}
			abs, err := a.FS.Abs(subdirFullPath)
			if err != nil {
				a.Logger.Println("Cannot compute abs path for:", childPath)
				f.Close()
				continue
			}
			nextEntry := &model.DirectoryEntry{File: f, AbsPath: abs}
			stack = append(stack, &RecursiveLookupFrame{entry: nextEntry, header: childPath})
		}
	}
}

// ProcessFiles converts a slice of file entries into a *model.Directory for printing.
func (a *App) ProcessFiles(files []model.FileEntry) *model.Directory {
	t := new(model.Directory)
	isLong := a.Config.LongListingMode != model.LongListingNone

	for _, fileEntry := range files {
		entry := a.buildEntry(fileEntry.AbsPath, fileEntry.Info, isLong)
		t.Files = append(t.Files, entry)
	}

	return t
}

// ProcessDirectory reads the contents of the given directory, builds a *model.Directory
// that contains *model.Entry objects for each item, and returns it.
func (a *App) ProcessDirectory(d *model.DirectoryEntry) (*model.Directory, error) {
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

// Reads directory contents, creates *model.Entry objects, adds special entries (.), (..)
// and sets up Git statuses if requested.
func (a *App) populateDirectory(d *model.DirectoryEntry, dirStat fs.FileInfo) (*model.Directory, error) {
	t := new(model.Directory)
	isLong := a.Config.LongListingMode != model.LongListingNone

	// If we need to show the current directory as an entry
	if a.Config.AllMode == model.IncludeAll || a.Config.Directory {
		t.Info = a.buildEntry(d.Name(), dirStat, isLong)

		if !a.Config.Directory {
			t.Info.Name = "."
			t.Info.Ext = ""
		}

		if !a.Config.DisableIcon {
			t.Info.Icon = format.GetOpenDirIcon()
		}
	}

	if a.Config.Directory {
		t.Files = append(t.Files, t.Info)
		return t, nil
	}

	entries, err := d.File.ReadDir(0)
	// We proceed even if err != nil, as entries may contain a partial list.

	// If Git status is requested, prepare the repository info map.
	var gitRepoStatus map[string]string
	if a.Config.GitStatus {
		gitRepoStatus = a.gitStatusFor(d.Name())
	}

	showHidden := a.Config.AllMode != model.IncludeDefault

	// Build entries for each file
	for _, de := range entries {
		name := de.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		fullpath := a.FS.Join(d.Name(), name)
		fi, infoErr := de.Info() // fi might be nil on error
		if infoErr != nil {
			a.Logger.Printf("cannot access %q: %v\n", fullpath, infoErr)
			a.ExitCode.SetMinor()
		}

		entry := a.buildEntry(fullpath, fi, isLong)

		// If we couldn't get full info but we have type from DirEntry, fill
		// in the indicator and pick a fallback icon. The inspector handles
		// symlink-icon resolution for entries that do have FileInfo.
		if fi == nil {
			entry.Mode = "???????????"
			if de.IsDir() {
				entry.Indicator = "/"
			} else if de.Type()&iofs.ModeSymlink != 0 {
				entry.Indicator = "@"
			}
			if !a.Config.DisableIcon {
				entry.Icon = format.GetIcon(entry.Name, entry.Ext, entry.Indicator)
			}
		}

		// If Git status is available, attach it.
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

	if a.Config.AllMode == model.IncludeAll {
		if t.Info != nil {
			t.Files = append(t.Files, t.Info)
		}

		pp := a.FS.Dir(d.Name())
		pStat, _ := a.FS.Lstat(pp)

		parentEntry := a.buildEntry(pp, pStat, isLong)
		parentEntry.Name = ".."
		parentEntry.Ext = ""

		if !a.Config.DisableIcon {
			parentEntry.Icon = format.GetOpenDirIcon()
		}

		t.Files = append(t.Files, parentEntry)
		t.Parent = parentEntry
	}
	return t, err
}

// inspectorFor returns an Inspector configured to populate exactly the fields
// the current request needs.
func (a *App) inspectorFor(isLong bool) *inspect.Inspector {
	showOwner := a.Config.LongListingMode == model.LongListingDefault ||
		a.Config.LongListingMode == model.LongListingOwner
	showGroup := !a.Config.NoGroup &&
		(a.Config.LongListingMode == model.LongListingDefault ||
			a.Config.LongListingMode == model.LongListingGroup)
	return inspect.New(a.FS, inspect.DefaultIconResolver(), inspect.Options{
		Long:            isLong,
		ShowOwner:       showOwner,
		ShowGroup:       showGroup,
		ShowInode:       a.Config.ShowInodeNumber,
		ShowBlocks:      a.Config.ShowBlockSize,
		ResolveSymlinks: !a.Config.DisableIcon,
		DisableIcon:     a.Config.DisableIcon,
	})
}

// Constructs a *model.Entry from a given path, fs.FileInfo, and whether we are
// in a long-listing context.
func (a *App) buildEntry(fullPath string, fi fs.FileInfo, isLong bool) *model.Entry {
	if fi == nil {
		entry := &model.Entry{}
		entry.Name = a.FS.Base(fullPath)
		entry.Ext = a.FS.Ext(entry.Name)
		entry.Name = entry.Name[0 : len(entry.Name)-len(entry.Ext)]
		entry.Mode = "???????????"
		entry.Owner = "?"
		entry.Group = "?"
		return entry
	}

	insp := a.inspectorFor(isLong)
	ie := insp.Inspect(fullPath, fi)

	modeStr := ""
	owner := ie.Owner
	group := ie.Group
	if isLong {
		modeStr = inspect.ModeString(ie.Mode, ie.Sticky, ie.StickyX, ie.HasXAttr)
		// Legacy renderer expects the group column pre-padded with " %v  "
		// (formatting baked into the old ctw column widths); preserve that
		// shape here until the renderer consumes InspectedEntry directly.
		if group != "" {
			group = fmt.Sprintf(" %v  ", group)
		}
	}

	entry := inspect.ToLegacy(ie, modeStr, owner, group, a.FS)
	return entry
}

// PrintDirectory sorts the directory's files according to the app config and prints them.
func (a *App) PrintDirectory(d *model.Directory) {
	if d == nil {
		return
	}

	format.SetLessFunction(d, a.Config.SortMode)
	d.Sort(a.Config.SortMode, a.Config.Reverse)

	render.Render(a.Writer, d.Files, render.Options{
		Mode:          a.renderMode(),
		ShowIcon:      !a.Config.DisableIcon,
		ShowInode:     a.Config.ShowInodeNumber,
		ShowBlocks:    a.Config.ShowBlockSize,
		HumanReadable: a.Config.HumanReadable,
		TimeFormatter: entryTimeFormatter{tf: a.Config.TimeFormatter},
	})
}

func (a *App) renderMode() render.Mode {
	switch {
	case a.Config.LongListingMode != model.LongListingNone:
		return render.ModeLong
	case a.Config.OneFilePerLine:
		return render.ModeOneFilePerLine
	default:
		return render.ModeShort
	}
}

// entryTimeFormatter adapts the app's *time.Time-based formatter to the
// renderer's per-entry interface.
type entryTimeFormatter struct {
	tf format.Timestamp
}

func (e entryTimeFormatter) Format(m *model.Entry) string {
	return e.tf.Format(&m.ModTime)
}
