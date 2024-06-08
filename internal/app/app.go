package app

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/canta2899/logo-ls/assets"
	"github.com/canta2899/logo-ls/internal/ctw"
	"github.com/canta2899/logo-ls/internal/format"
	"github.com/canta2899/logo-ls/internal/git_utils"
	"github.com/canta2899/logo-ls/internal/model"
)

type App struct {
	Config        *Config
	Writer        io.Writer
	ExitCode      model.ExitCode
	TerminalWidth int
	Logger        *log.Logger
}

type Args struct {
	Files []model.FileInfo
	Dirs  []model.File
}

func (a *App) Exit() {
	os.Exit(int(a.ExitCode))
}

func (a *App) Write(buf *bytes.Buffer) {
	_, err := io.Copy(a.Writer, buf)

	if err != nil {
		panic(err)
	}
}

func (a *App) GetArguments() *Args {
	dirs := a.Config.FileList

	sort.Slice(dirs, func(i, j int) bool {
		return dirs[i].Path < dirs[j].Path
	})

	args := &Args{}

	// segregate args in files and dirs, and print error for those which cannot be opened
	for _, v := range dirs {

		d, err := os.Open(v.AbsolutePath)

		if err != nil {
			log.Printf("cannot access %q: %v\n", v, err)
			d.Close()
			a.ExitCode.SetSerious()
			continue
		}

		ds, err := d.Stat()
		if err != nil {
			log.Printf("cannot access %q: %v\n", v, err)
			d.Close()
			a.ExitCode.SetSerious()
			continue
		}

		if ds.IsDir() {
			args.Dirs = append(args.Dirs, model.File{
				File:    *d,
				AbsPath: v.AbsolutePath,
			})
		} else {
			args.Files = append(args.Files, model.FileInfo{
				FileInfo: ds,
				AbsPath:  v.AbsolutePath,
			})
		}
	}

	return args
}

func (a *App) GetCtw() ctw.CTW {
	var out ctw.CTW

	if a.Config.LongListingMode != model.LongListingNone {
		out = ctw.NewLongCTW(9)
	} else if a.Config.OneFilePerLine {
		out = ctw.NewLongCTW(4)
	} else {
		out = ctw.NewStandardCTW(a.TerminalWidth)
	}

	out.DisplayColor(!a.Config.DisableColor)

	if a.Config.DisableColor {
		model.OpenDirIcon = assets.Icon_Def["diropen"].GetGlyph() + " "
	}

	if a.Config.DisableIcon {
		model.OpenDirIcon = ""
	}

	return out
}

func (a *App) blockSize(block int64) string {
	if a.Config.ShowBlockSize {
		return a.getSizeFromFormat(block)
	}

	return ""
}

func (a *App) Print(d *model.Dir) {

	format.SetLessFunction(d, a.Config.SortMode)
	d.Sort(a.Config.SortMode, a.Config.Reverse)

	buf := bytes.NewBuffer([]byte(""))
	lineCtw := a.GetCtw()

	switch {

	case a.Config.LongListingMode != model.LongListingNone:
		for _, v := range d.Files {
			lineCtw.AddRow(
				a.blockSize(v.Blocks),
				v.Mode,
				v.Owner,
				v.Group,
				a.getSizeFromFormat(v.Size),
				v.ModTime.Format(a.Config.TimeFormat),
				v.Icon,
				v.Name+v.Ext+v.Indicator,
				v.GitStatus)

			lineCtw.IconColor(v.IconColor)
		}

	case a.Config.OneFilePerLine:
		for _, v := range d.Files {
			lineCtw.AddRow(a.blockSize(v.Blocks), v.Icon, v.Name+v.Ext+v.Indicator, v.GitStatus)
			lineCtw.IconColor(v.IconColor)
		}

	default:
		for _, v := range d.Files {
			lineCtw.AddRow(a.blockSize(v.Blocks), v.Icon, v.Name+v.Ext+v.Indicator, v.GitStatus)
			lineCtw.IconColor(v.IconColor)
		}
	}

	lineCtw.Flush(buf)
	a.Write(buf)
}

func (a *App) ListFiles(files []model.FileInfo) *model.Dir {

	isLong := a.Config.LongListingMode != model.LongListingNone

	t := new(model.Dir)

	for _, v := range files {
		name := v.Name()
		f := new(model.Entry)
		f.Ext = filepath.Ext(name)
		f.Name = name[0 : len(name)-len(f.Ext)]
		f.Indicator = format.GetIndicator(v.AbsPath, isLong)
		f.Size = v.Size()
		f.ModTime = v.ModTime()

		if isLong {
			f.Mode = v.Mode().String()
			f.ModeBits = uint32(v.Mode())
			f.Owner, f.Group = model.GetOwnerGroupInfo(v, a.Config.NoGroup, a.Config.LongListingMode)
		}

		if a.Config.ShowBlockSize {
			model.DirBlocks(f, v)
		}

		if !a.Config.DisableIcon {
			f.Icon, f.IconColor = format.GetIcon(f.Name, f.Ext, f.Indicator)
			if a.Config.DisableColor {
				f.IconColor = ""
			}
		}

		t.Files = append(t.Files, f)
	}
	return t
}

func (a *App) ListDirs(d *model.File) (*model.Dir, error) {
	// some flag variable combinations
	long := a.Config.LongListingMode != model.LongListingNone
	currentDir := a.Config.AllMode == model.IncludeAll || a.Config.Directory
	showHidden := a.Config.AllMode != model.IncludeDefault

	t := new(model.Dir)

	// filing current dir info
	t.Info = new(model.Entry)
	t.Info.Name = "."
	ds, err := d.Stat()
	if err != nil {
		return nil, err
	}

	// getting Git Status of the entire repository
	var gitRepoStatus map[string]string // could be nil
	if a.Config.GitStatus {
		gitRepoStatus = git_utils.GetFilesGitStatus(d.Name()) // returns map or nil
		if len(gitRepoStatus) == 0 {
			gitRepoStatus = nil
		}
	}

	if currentDir {
		t.Info.Size = ds.Size()
		t.Info.ModTime = ds.ModTime()
		if long {
			t.Info.Mode = ds.Mode().String()
			t.Info.ModeBits = uint32(ds.Mode())
			t.Info.Owner, t.Info.Group = model.GetOwnerGroupInfo(ds, a.Config.NoGroup, a.Config.LongListingMode)
		}
		if a.Config.ShowBlockSize {
			model.DirBlocks(t.Info, ds)
		}
		if !a.Config.DisableIcon {
			t.Info.Icon = assets.Icon_Def["diropen"].GetGlyph()
			if !a.Config.DisableColor {
				t.Info.IconColor = assets.Icon_Def["diropen"].GetColor(1)
			}
		}
	}

	// don't fill files info if the -d flag is passed
	if a.Config.Directory {
		t.Files = append(t.Files, t.Info)
		return t, nil
	}

	files, err := d.Readdir(0)
	for _, v := range files {
		name := v.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}

		fullpath := filepath.Join(d.Name(), name)
		f := new(model.Entry)
		f.Ext = filepath.Ext(name)
		f.Name = name[0 : len(name)-len(f.Ext)]
		f.Indicator = format.GetIndicator(fullpath, long)
		f.Size = v.Size()
		f.ModTime = v.ModTime()
		if long {
			f.Mode = v.Mode().String()
			f.ModeBits = uint32(v.Mode())
			f.Owner, f.Group = model.GetOwnerGroupInfo(v, a.Config.NoGroup, a.Config.LongListingMode)
		}
		if a.Config.ShowBlockSize {
			model.DirBlocks(f, v)
		}

		if !a.Config.DisableIcon {
			f.Icon, f.IconColor = format.GetIcon(f.Name, f.Ext, f.Indicator)
			if format.IsLink(fullpath) {
				if s, err := filepath.EvalSymlinks(fullpath); err == nil {
					realExt := filepath.Ext(s)
					realName := s[0 : len(s)-len(realExt)]
					realIndicator := format.GetIndicator(s, long)
					f.Icon, f.IconColor = format.GetIcon(realName, realExt, realIndicator)
				}
			}
			if a.Config.DisableColor {
				f.IconColor = ""
			}
		}

		if gitRepoStatus != nil {
			if v.IsDir() {
				f.GitStatus = gitRepoStatus[model.PathSeparator+v.Name()+model.PathSeparator]
			} else {
				f.GitStatus = gitRepoStatus[model.PathSeparator+v.Name()]
			}
		}

		t.Files = append(t.Files, f)
		if v.IsDir() {
			t.Dirs = append(t.Dirs, name+"/")
		}
	}

	// if -a flag is passed then only eval parent dir and append to files
	if a.Config.AllMode == model.IncludeAll {
		t.Files = append(t.Files, t.Info)
		p, err := filepath.Abs(d.Name())
		if err != nil {
			// partial *dir (without parent dir) and error
			return t, err
		}
		pp := filepath.Dir(p)
		pds, err := os.Lstat(pp)
		if err != nil {
			// partial *dir (without parent dir) and error
			return t, err
		}
		t.Parent = new(model.Entry)
		t.Parent.Name = ".."
		t.Parent.Size = pds.Size()
		t.Parent.ModTime = pds.ModTime()
		if long {
			t.Parent.Mode = pds.Mode().String()
			t.Parent.ModeBits = uint32(pds.Mode())
			t.Parent.Owner, t.Parent.Group = model.GetOwnerGroupInfo(pds, a.Config.NoGroup, a.Config.LongListingMode)
		}
		if a.Config.ShowBlockSize {
			model.DirBlocks(t.Parent, pds)
		}
		if !a.Config.DisableIcon {
			t.Parent.Icon = assets.Icon_Def["diropen"].GetGlyph()
			if !a.Config.DisableColor {
				t.Parent.IconColor = assets.Icon_Def["diropen"].GetColor(1)
			}
		}
		t.Files = append(t.Files, t.Parent)
	}

	// return *dir with no error
	// or partial *dir with error (produced by Readdir)
	return t, err
}

func (a *App) RecurseDirs(d *model.File, startingAbsolutePath string) {
	dd, err := a.ListDirs(d)
	d.Close()
	if err != nil {
		log.Printf("partial access to %q: %v\n", d.Name(), err)
		// sysState.ExitCode(sysState.Code_Minor)
	}

	a.Print(dd)

	if len(dd.Dirs) == 0 {
		return
	}
	// at this point dd.print has sorted the children files
	// but not using it instead printing children in directory order
	temp := make([]string, len(dd.Dirs))
	sort.Strings(dd.Dirs)
	for i, v := range dd.Dirs {
		rel, err := filepath.Rel(startingAbsolutePath, d.Name())

		if err == nil {
			temp[i] = filepath.Join(rel, v)
		} else {
			temp[i] = filepath.Join(d.Name(), v)
		}
	}
	for _, v := range temp {
		fmt.Printf("\n%s:\n", model.OpenDirIcon+v)
		f, err := os.Open(v)
		if err != nil {
			log.Printf("cannot access %q: %v\n", v, err)
			f.Close()
			// sysState.ExitCode(sysState.Code_Minor)
			continue
		}
		abs, err := filepath.Abs(v)

		if err != nil {
			log.Println("Cannot compute abs path")
			f.Close()
			continue
		}

		next := &model.File{File: *f, AbsPath: abs}
		a.RecurseDirs(next, startingAbsolutePath)
	}
}

func (a *App) Run() {

	args := a.GetArguments()

	// process and display all files
	if len(args.Files) > 0 {
		a.Print(a.ListFiles(args.Files))
		if len(args.Dirs) > 0 {
			fmt.Println()
		}
	}

	// process and display all the dirs in arg
	if a.Config.Recursive {
		// use recursive func
		for i, v := range args.Dirs {
			if i > 0 {
				fmt.Println()
			}

			currentAbsolutePath, _ := filepath.Abs(".")
			fileRelativePath, err := filepath.Rel(currentAbsolutePath, v.Name())

			if err == nil {
				fmt.Printf("%s:\n", model.OpenDirIcon+fileRelativePath)
			} else {
				fmt.Printf("%s:\n", model.OpenDirIcon+v.Name())
			}

			if a.Config.GitStatus {
				git_utils.GitRepoCompute()
			}

			a.RecurseDirs(&v, currentAbsolutePath)
		}
	} else {
		pName := len(args.Dirs) > 1
		for i, v := range args.Dirs {
			if pName {
				fmt.Printf("%s:\n", model.OpenDirIcon+v.Name())
			}
			if a.Config.GitStatus {
				git_utils.GitRepoCompute()
			}
			d, err := a.ListDirs(&v)
			v.Close()
			if err != nil {
				log.Printf("partial access to %q: %v\n", v.Name(), err)
				a.ExitCode.SetSerious()
			}
			a.Print(d)
			if i < len(args.Dirs)-1 {
				fmt.Println()
			}
		}
	}
}

func (a *App) getSizeFromFormat(b int64) string {

	if !a.Config.HumanReadable {
		return fmt.Sprintf("%d", b)
	}

	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d", b)
	}

	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c",
		float64(b)/float64(div), "KMGTPE"[exp])
}
