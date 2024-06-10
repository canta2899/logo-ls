package icons_test

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"testing"

	"github.com/canta2899/logo-ls/ctw"
	"github.com/canta2899/logo-ls/icons"
	"golang.org/x/term"
)

func TestFileIcons(t *testing.T) {
	log.Println("Printing each supported file name and ext by the icon pack")

	// get terminal width
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
			i := icons.IconSet[v]
			fmt.Fprintln(os.Stderr)
			buf := bytes.NewBuffer([]byte(""))
			log.Println("Printing files of type", i.GetColor(1)+v+"\033[0m")
			w := ctw.NewStandardCTW(terminalWidth)
			for f, d := range icons.IconFileName {
				if d == i {
					w.AddRow("    ", d.GetGlyph(), f, "")
					w.IconColor(d.GetColor(1))
				}
			}
			w.Flush(buf)
			io.Copy(os.Stderr, buf)

			buf = bytes.NewBuffer([]byte(""))
			log.Println("Printing extentions of type", i.GetColor(1)+v+"\033[0m")
			w = ctw.NewStandardCTW(terminalWidth)
			for e, d := range icons.IconExt {
				if d == i {
					w.AddRow("    ", d.GetGlyph(), e, "")
					w.IconColor(d.GetColor(1))
				}
			}
			w.Flush(buf)
			io.Copy(os.Stderr, buf)
		})
	}
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
				w.AddRow("    ", set[v].GetGlyph(), v, "")
				w.IconColor(set[v].GetColor(1))
			}
			w.Flush(buf)
			io.Copy(os.Stdout, buf)
			fmt.Fprintln(os.Stdout)
		})
	}
}
