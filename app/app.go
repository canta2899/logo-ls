package app

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/canta2899/logo-ls/ctw"
	"github.com/canta2899/logo-ls/format"
	"github.com/canta2899/logo-ls/icons"
	"github.com/canta2899/logo-ls/model"
)

// Represents the main application that holds configuration, a writer, exit codes, and a logger.
type App struct {
	Config   *Config
	Writer   io.Writer
	ExitCode model.ExitCode
	Logger   *log.Logger
}

// Stores the parsed command-line arguments as separate files and directories.
type Args struct {
	Files []model.FileEntry
	Dirs  []model.DirectoryEntry
}

// Terminates the program with the current stored exit code.
func (a *App) Exit() {
	os.Exit(int(a.ExitCode))
}

// Copies a buffer to the App's writer, panicking on error.
func (a *App) Write(buf *bytes.Buffer) {
	if _, err := io.Copy(a.Writer, buf); err != nil {
		panic(err)
	}
}

// Parses the configured file list and categorizes them as files or directories.
// It also sets an error exit code for entries that are not accessible.
func (a *App) GetArguments() *Args {
	// Sort all user inputs.
	sort.Slice(a.Config.FileList, func(i, j int) bool {
		return a.Config.FileList[i] < a.Config.FileList[j]
	})

	args := &Args{}

	for _, argPath := range a.Config.FileList {
		abs, err := filepath.Abs(argPath)
		if err != nil {
			a.Logger.Printf("cannot get absolute path for %q: %v\n", argPath, err)
			a.ExitCode.SetSerious()
			continue
		}

		f, err := os.Open(abs)
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
				File:    *f,
				AbsPath: abs,
			})
		} else {
			args.Files = append(args.Files, model.FileEntry{
				FileInfo: fi,
				AbsPath:  abs,
			})
			f.Close() // no need to keep file handles for single files
		}
	}
	return args
}

// Main entry point that orchestrates listing files/directories and printing results.
func (a *App) Run() {
	args := a.GetArguments()

	// Process and display all files first.
	if len(args.Files) > 0 {
		filesDir := a.ProcessFiles(args.Files)
		a.PrintDirectory(filesDir)
		if len(args.Dirs) > 0 {
			fmt.Println()
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
	currentAbs, _ := filepath.Abs(".")
	openDirIcon := model.OpenDirIcon
	if a.Config.DisableIcon {
		openDirIcon = ""
	}

	for i, dirEntry := range dirs {
		if i > 0 {
			fmt.Println()
		}

		relName := dirEntry.Name()
		if rel, err := filepath.Rel(currentAbs, dirEntry.Name()); err == nil {
			relName = rel
		}

		fmt.Printf("%s:\n", openDirIcon+relName)

		a.recurseDirectory(&dirEntry, currentAbs)
	}
}

// Prints each directory (but not subdirectories).
func (a *App) processDirsNonRecursively(dirs []model.DirectoryEntry) {
	pName := len(dirs) > 1
	openDirIcon := model.OpenDirIcon
	if a.Config.DisableIcon {
		openDirIcon = ""
	}

	for i, dirEntry := range dirs {
		if pName {
			fmt.Printf("%s:\n", openDirIcon+dirEntry.Name())
		}

		d, err := a.ProcessDirectory(&dirEntry)
		dirEntry.Close()
		if err != nil {
			a.Logger.Printf("partial access to %q: %v\n", dirEntry.Name(), err)
			a.ExitCode.SetSerious()
		}

		a.PrintDirectory(d)
		if i < len(dirs)-1 {
			fmt.Println()
		}
	}
}

// Processes a directory, prints it, and recurses through subdirectories if any.
func (a *App) recurseDirectory(dir *model.DirectoryEntry, startingAbsolutePath string) {
	d, err := a.ProcessDirectory(dir)
	dir.Close()
	if err != nil {
		a.Logger.Printf("partial access to %q: %v\n", dir.Name(), err)
		a.ExitCode.SetMinor()
	}

	a.PrintDirectory(d)

	if len(d.Dirs) == 0 {
		return
	}

	sort.Strings(d.Dirs)
	for _, subdir := range d.Dirs {
		childPath := filepath.Join(dir.Name(), subdir)
		if rel, err := filepath.Rel(startingAbsolutePath, childPath); err == nil {
			childPath = rel
		}

		fmt.Printf("\n%s:\n", model.OpenDirIcon+childPath)

		f, err := os.Open(filepath.Join(dir.Name(), subdir))
		if err != nil {
			a.Logger.Printf("cannot access %q: %v\n", childPath, err)
			a.ExitCode.SetMinor()
			continue
		}
		abs, err := filepath.Abs(filepath.Join(dir.Name(), subdir))
		if err != nil {
			a.Logger.Println("Cannot compute abs path for:", childPath)
			f.Close()
			continue
		}
		next := &model.DirectoryEntry{File: *f, AbsPath: abs}
		a.recurseDirectory(next, startingAbsolutePath)
	}
}

// Converts a slice of file entries into a *model.Directory for printing.
func (a *App) ProcessFiles(files []model.FileEntry) *model.Directory {
	t := new(model.Directory)
	isLong := a.Config.LongListingMode != model.LongListingNone

	for _, fileEntry := range files {
		entry := a.buildEntry(fileEntry.AbsPath, fileEntry.FileInfo, isLong)
		t.Files = append(t.Files, entry)
	}

	return t
}

// Reads the contents of the given directory, builds a *model.Directory
// that contains *model.Entry objects for each item, and returns it.
func (a *App) ProcessDirectory(d *model.DirectoryEntry) (*model.Directory, error) {
	defer func() {
		// Defer the close in case 'Readdir' triggers partial reads.
		_ = d.Close()
	}()

	dirStat, err := d.Stat()
	if err != nil {
		return nil, err
	}

	dirModel, err := a.populateDirectory(d, dirStat)
	return dirModel, err
}

// Reads directory contents, creates *model.Entry objects, adds special entries (.), (..)
// and sets up Git statuses if requested.
func (a *App) populateDirectory(d *model.DirectoryEntry, dirStat os.FileInfo) (*model.Directory, error) {
	t := new(model.Directory)
	isLong := a.Config.LongListingMode != model.LongListingNone

	// If we need to show the current directory as an entry
	if a.Config.AllMode == model.IncludeAll || a.Config.Directory {
		t.Info = a.buildEntry(d.Name(), dirStat, isLong)

		if pwd, err := os.Getwd(); err == nil {
			if d.Name() == pwd {
				t.Info.Name = "."
			}
		}

		if !a.Config.DisableIcon {
			t.Info.Icon = icons.IconDef["diropen"].GetGlyph()
			t.Info.IconColor = icons.IconDef["diropen"].GetColor()
		}
	}

	if a.Config.Directory {
		t.Files = append(t.Files, t.Info)
		return t, nil
	}

	files, err := d.Readdir(0)
	if err != nil {
		return t, err
	}

	// If Git status is requested, prepare the repository info map.
	var gitRepoStatus map[string]string
	if a.Config.GitStatus {
		gitRepoStatus = format.GetFilesGitStatus(d.Name())
		if len(gitRepoStatus) == 0 {
			gitRepoStatus = nil
		}
	}

	showHidden := a.Config.AllMode != model.IncludeDefault

	// Build entries for each file
	for _, fi := range files {
		name := fi.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		fullpath := filepath.Join(d.Name(), name)
		entry := a.buildEntry(fullpath, fi, isLong)

		if !a.Config.DisableIcon {
			if format.IsLink(fullpath) {
				if s, err := filepath.EvalSymlinks(fullpath); err == nil {
					realExt := filepath.Ext(s)
					realName := s[0 : len(s)-len(realExt)]
					realIndicator := format.GetIndicator(s, isLong)
					entry.Icon, entry.IconColor = format.GetIcon(realName, realExt, realIndicator)
				}
			}
		}

		// If Git status is available, attach it.
		if gitRepoStatus != nil {
			entry.GitStatus = gitRepoStatus[fi.Name()+model.PathSeparator]
			if entry.GitStatus == "" {
				entry.GitStatus = gitRepoStatus[fi.Name()]
			}
		}

		t.Files = append(t.Files, entry)
		if fi.IsDir() {
			t.Dirs = append(t.Dirs, name+"/")
		}
	}

	if a.Config.AllMode == model.IncludeAll {
		if t.Info != nil {
			t.Files = append(t.Files, t.Info)
		}

		pp := filepath.Dir(d.Name())
		pStat, err2 := os.Lstat(pp)

		if err2 == nil { // if we can't stat parent, skip adding
			parentEntry := a.buildEntry(pp, pStat, isLong)
			parentEntry.Name = ".."

			// Overwrite icon for parent
			if !a.Config.DisableIcon {
				parentEntry.Icon = icons.IconDef["diropen"].GetGlyph()
				parentEntry.IconColor = icons.IconDef["diropen"].GetColor()
			}

			t.Files = append(t.Files, parentEntry)
			t.Parent = parentEntry
		}
	}
	return t, err
}

// Constructs a *model.Entry from a given path, os.FileInfo, and whether we are
// in a long-listing context.
func (a *App) buildEntry(fullPath string, fi os.FileInfo, isLong bool) *model.Entry {
	entry := &model.Entry{}

	name := fi.Name()
	entry.Ext = filepath.Ext(name)
	entry.Name = name[0 : len(name)-len(entry.Ext)]
	entry.Size = fi.Size()
	entry.ModTime = fi.ModTime()
	entry.Indicator = format.GetIndicator(fullPath, isLong)

	if a.Config.ShowInodeNumber {
		entry.InodeNumber = format.GetInodeNumber(fullPath)
	}

	if isLong {
		entry.Mode = format.GetModeExtended(&fi, fullPath)
		entry.ModeBits = uint32(fi.Mode())
		entry.NumHardLinks = format.GetHardLinkCount(fullPath)
		owner, group := model.GetOwnerGroupInfo(fi, a.Config.NoGroup, a.Config.LongListingMode)
		entry.Owner, entry.Group = owner, group
	}

	if a.Config.ShowBlockSize {
		model.DirBlocks(entry, fi)
	}

	if !a.Config.DisableIcon {
		entry.Icon, entry.IconColor = format.GetIcon(entry.Name, entry.Ext, entry.Indicator)
	}

	return entry
}

// Sorts the directory's files according to the app config and prints them.
func (a *App) PrintDirectory(d *model.Directory) {
	if d == nil {
		return
	}

	format.SetLessFunction(d, a.Config.SortMode)
	d.Sort(a.Config.SortMode, a.Config.Reverse)

	lineCtw := a.getCTW()

	switch {
	case a.Config.LongListingMode != model.LongListingNone:
		for _, f := range d.Files {
			lineCtw.AddRow(
				f.IconColor,
				a.blockSizeWithInode(f),
				f.Mode,
				strconv.Itoa(int(f.NumHardLinks)),
				f.Owner,
				f.Group,
				format.GetFormattedSize(f.Size, a.Config.HumanReadable),
				a.Config.TimeFormatter.Format(&f.ModTime),
				f.Icon,
				f.Name+f.Ext+f.Indicator,
				f.GitStatus,
			)
		}

	case a.Config.OneFilePerLine:
		for _, f := range d.Files {
			lineCtw.AddRow(
				f.IconColor,
				a.blockSizeWithInode(f),
				f.Icon,
				f.Name+f.Ext+f.Indicator,
				f.GitStatus,
			)
		}

	default:
		for _, f := range d.Files {
			lineCtw.AddRow(
				f.IconColor,
				a.blockSizeWithInode(f),
				f.Icon,
				f.Name+f.Ext+f.Indicator,
				f.GitStatus,
			)
		}
	}

	// Flush to buffer and write out
	buf := new(bytes.Buffer)
	lineCtw.Flush(buf)
	a.Write(buf)
}

// Instantiates the proper ctw.CTW based on app config.
func (a *App) getCTW() ctw.CTW {
	longMode := a.Config.LongListingMode != model.LongListingNone
	return ctw.NewCTW(longMode, a.Config.OneFilePerLine, !a.Config.DisableIcon)
}

// Generates the block size and optional inode as a single string.
func (a *App) blockSizeWithInode(e *model.Entry) string {
	var parts []string

	if a.Config.ShowInodeNumber {
		parts = append(parts, e.InodeNumber)
	}

	if a.Config.ShowBlockSize {
		parts = append(parts, format.GetFormattedSize(e.Blocks, a.Config.HumanReadable))
	}

	return strings.TrimSpace(strings.Join(parts, " "))
}
