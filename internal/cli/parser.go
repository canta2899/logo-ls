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

		switch {
		case strings.HasPrefix(arg, "--"):
			consumed, err := p.parseLong(args, i)
			if err != nil {
				return err
			}
			i += consumed
		case strings.HasPrefix(arg, "-") && len(arg) > 1:
			if err := p.parseShortCluster(arg[1:]); err != nil {
				return err
			}
		default:
			p.Args = append(p.Args, arg)
		}
	}
	return nil
}

// parseLong handles a "--name[=value]" token. Returns the number of additional
// args consumed (0 or 1 for a value following the flag).
func (p *Parser) parseLong(args []string, i int) (int, error) {
	name := args[i][2:]
	var inlineValue string
	hasInline := false
	if eq := strings.IndexByte(name, '='); eq >= 0 {
		inlineValue = name[eq+1:]
		name = name[:eq]
		hasInline = true
	}

	f := p.findLong(name)
	if f == nil {
		return 0, fmt.Errorf("unknown flag: --%s", name)
	}
	return p.assignLong(f, name, args, i, inlineValue, hasInline)
}

func (p *Parser) findLong(name string) *Flag {
	for _, f := range p.Flags {
		if f.Long == name {
			return f
		}
	}
	return nil
}

func (p *Parser) assignLong(f *Flag, name string, args []string, i int, inlineValue string, hasInline bool) (int, error) {
	switch f.Type {
	case FlagBool:
		if hasInline {
			return 0, fmt.Errorf("flag --%s does not take a value", name)
		}
		*(f.Value.(*bool)) = true
		return 0, nil
	case FlagString:
		if hasInline {
			*(f.Value.(*string)) = inlineValue
			return 0, nil
		}
		if i+1 >= len(args) {
			return 0, fmt.Errorf("flag --%s requires a value", name)
		}
		*(f.Value.(*string)) = args[i+1]
		return 1, nil
	}
	return 0, nil
}

func (p *Parser) parseShortCluster(chars string) error {
	for _, char := range chars {
		if !p.setShortBool(char) {
			return fmt.Errorf("unknown flag: -%c", char)
		}
	}
	return nil
}

func (p *Parser) setShortBool(char rune) bool {
	for _, f := range p.Flags {
		if f.Short == char && f.Type == FlagBool {
			*(f.Value.(*bool)) = true
			return true
		}
	}
	return false
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
