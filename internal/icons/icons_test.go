package icons_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/canta2899/logo-ls/internal/render/columns"
	"github.com/canta2899/logo-ls/internal/icons"
	"golang.org/x/term"
)

func TestFileIcons(t *testing.T) {
	log.Println("Printing each supported file name and ext by the icon pack")

	terminalWidth, _, e := term.GetSize(int(os.Stdout.Fd()))
	if e != nil {
		terminalWidth = 80
	}

	ks := make([]string, 0)
	for k := range icons.IconSet {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	for _, v := range ks {
		t.Run("Testing icon: "+v, func(st *testing.T) {
			runIconSubtest(v, terminalWidth)
		})
	}
}

func runIconSubtest(name string, terminalWidth int) {
	i := icons.IconSet[name]
	fmt.Fprintln(os.Stderr)
	log.Println("Printing files of type", i.GetColor()+name+"\033[0m")
	writeMatchingEntries(terminalWidth, i, icons.IconFileName)

	log.Println("Printing extentions of type", i.GetColor()+name+"\033[0m")
	writeMatchingEntries(terminalWidth, i, icons.IconExt)
}

func writeMatchingEntries(terminalWidth int, target *icons.IconInfo, set map[string]*icons.IconInfo) {
	buf := bytes.NewBuffer(nil)
	w := ctw.NewStandardCTW(terminalWidth)
	for name, info := range set {
		if info == target {
			w.AddRow(info.GetColor(), "    ", info.GetGlyph(), name, "")
		}
	}
	w.Flush(buf)
	io.Copy(os.Stderr, buf)
}

func TestIconDisplay(t *testing.T) {
	// get terminal width
	terminalWidth, _, e := term.GetSize(int(os.Stdout.Fd()))
	if e != nil {
		terminalWidth = 80
	}

	temp := [2]map[string]*icons.IconInfo{icons.IconSet, icons.IconDef}

	for i, set := range temp {
		t.Run(fmt.Sprintf("Icon Set %d", i+1), func(st *testing.T) {
			//sorting alphabetically
			ks := make([]string, 0)
			for k := range set {
				ks = append(ks, k)
			}
			sort.Strings(ks)

			// display icons
			buf := bytes.NewBuffer([]byte("\n"))
			w := ctw.NewStandardCTW(terminalWidth)
			for _, v := range ks {
				w.AddRow(set[v].GetColor(), "    ", set[v].GetGlyph(), v, "")
			}
			w.Flush(buf)
			io.Copy(os.Stdout, buf)
			fmt.Fprintln(os.Stdout)
		})
	}
}
