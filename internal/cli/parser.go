package cli

import (
	"fmt"
	"strings"
)

type FlagType int

const (
	FlagBool FlagType = iota
	FlagString
)

type Flag struct {
	Short       rune
	Long        string
	Description string
	Type        FlagType
	Value       any
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

func (p *Parser) String(long, defaultValue, description string) *string {
	val := defaultValue
	f := &Flag{
		Long:        long,
		Description: description,
		Type:        FlagString,
		Value:       &val,
	}
	p.Flags = append(p.Flags, f)
	return &val
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
			var inlineValue string
			hasInline := false
			if eq := strings.IndexByte(name, '='); eq >= 0 {
				inlineValue = name[eq+1:]
				name = name[:eq]
				hasInline = true
			}
			found := false
			for _, f := range p.Flags {
				if f.Long != name {
					continue
				}
				switch f.Type {
				case FlagBool:
					if hasInline {
						return fmt.Errorf("flag --%s does not take a value", name)
					}
					*(f.Value.(*bool)) = true
				case FlagString:
					if hasInline {
						*(f.Value.(*string)) = inlineValue
					} else {
						if i+1 >= len(args) {
							return fmt.Errorf("flag --%s requires a value", name)
						}
						i++
						*(f.Value.(*string)) = args[i]
					}
				}
				found = true
				break
			}
			if !found {
				return fmt.Errorf("unknown flag: --%s", name)
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
			if f.Type == FlagString {
				longStr += " <value>"
			}
		}
		fmt.Printf("  %-4s %-28s %s\n", shortStr, longStr, f.Description)
	}
}
