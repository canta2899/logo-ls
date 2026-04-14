package app

import (
	"fmt"
	"strings"
)

type FlagType int

const (
	FlagBool FlagType = iota
)

type Flag struct {
	Short       rune
	Long        string
	Description string
	Type        FlagType
	Value       interface{}
}

type Parser struct {
	Flags      []*Flag
	Args       []string
	Usage      string
	Parameters string
}

func NewParser() *Parser {
	return &Parser{
		Flags: []*Flag{},
		Args:  []string{},
	}
}

func (p *Parser) Bool(short rune, long string, description string) *bool {
	val := false
	f := &Flag{
		Short:       short,
		Long:        long,
		Description: description,
		Type:        FlagBool,
		Value:       &val,
	}
	p.Flags = append(p.Flags, f)
	return &val
}

func (p *Parser) Parse(args []string) error {
	// Skip executable name
	if len(args) > 0 {
		args = args[1:]
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if arg == "--" {
			p.Args = append(p.Args, args[i+1:]...)
			break
		}

		if strings.HasPrefix(arg, "--") {
			name := arg[2:]
			found := false
			for _, f := range p.Flags {
				if f.Long == name {
					if f.Type == FlagBool {
						*(f.Value.(*bool)) = true
						found = true
						break
					}
				}
			}
			if !found {
				return fmt.Errorf("unknown flag: %s", arg)
			}
		} else if strings.HasPrefix(arg, "-") && len(arg) > 1 {
			chars := arg[1:]
			for _, char := range chars {
				found := false
				for _, f := range p.Flags {
					if f.Short == char {
						if f.Type == FlagBool {
							*(f.Value.(*bool)) = true
							found = true
							break
						}
					}
				}
				if !found {
					return fmt.Errorf("unknown flag: -%c", char)
				}
			}
		} else {
			p.Args = append(p.Args, arg)
		}
	}
	return nil
}

func (p *Parser) PrintUsage() {
	fmt.Printf("Usage: logo-ls [options] %s\n", p.Parameters)
	fmt.Println("\nOptions:")
	for _, f := range p.Flags {
		shortStr := ""
		if f.Short != 0 {
			shortStr = fmt.Sprintf("-%c,", f.Short)
		}
		longStr := ""
		if f.Long != "" {
			longStr = fmt.Sprintf("--%s", f.Long)
		}
		fmt.Printf("  %-4s %-20s %s\n", shortStr, longStr, f.Description)
	}
}
