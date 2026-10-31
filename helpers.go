package cligo

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"
	"unicode"
)

func gatherArgs(args []string) []string {
	cmdArgs := []string{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			break
		}
		cmdArgs = append(cmdArgs, arg)
	}
	return cmdArgs
}

func gatherFlags(args []string) map[string]bool {
	flags := map[string]bool{}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			flags[arg] = false
		}
	}
	return flags
}

func getStringFlags(flagArgs []string, names []string) ([]string, error) {
	values := []string{}
	startIndex := 0
	for _, name := range names {
		for {
			matchingIndex := slices.Index(flagArgs[startIndex:], name)
			if matchingIndex == -1 {
				break
			}

			valIndex := startIndex + matchingIndex + 1
			if valIndex >= len(flagArgs) {
				return nil, fmt.Errorf("missing value for flag: %s", name)
			}
			values = append(values, flagArgs[valIndex])

			startIndex += matchingIndex + 2
			if startIndex >= len(flagArgs) {
				break
			}
		}
		startIndex = 0
	}

	return values, nil
}

func getBoolFlags(flagArgs []string, names []string) ([]bool, error) {
	values := []bool{}
	startIndex := 0
	for _, name := range names {
		for {
			matchingIndex := slices.Index(flagArgs[startIndex:], name)
			if matchingIndex == -1 {
				break
			}

			values = append(values, true)

			startIndex += matchingIndex + 1
			if startIndex >= len(flagArgs) {
				break
			}
		}
		startIndex = 0
	}

	return values, nil
}

func getIntFlags(valuesStr []string, name string) ([]int, error) {
	values := []int{}
	for _, v := range valuesStr {
		intVal, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value for flag %s: %s", name, v)
		}
		values = append(values, intVal)
	}
	return values, nil
}

func getUintFlags(valuesStr []string, name string) ([]uint, error) {
	intValues, err := getIntFlags(valuesStr, name)
	if err != nil {
		return nil, err
	}
	values := []uint{}
	for _, intVal := range intValues {
		if intVal < 0 {
			return nil, fmt.Errorf("invalid unsigned integer value for flag %s: %d", name, intVal)
		}
		values = append(values, uint(intVal))
	}
	return values, nil
}

func getFloat64Flags(valuesStr []string, name string) ([]float64, error) {
	values := []float64{}
	for _, v := range valuesStr {
		floatVal, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float64 value for flag %s: %s", name, v)
		}
		values = append(values, floatVal)
	}
	return values, nil
}

func getFloat32Flags(valuesStr []string, name string) ([]float32, error) {
	values := []float32{}
	for _, v := range valuesStr {
		floatVal, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid float32 value for flag %s: %s", name, v)
		}
		values = append(values, float32(floatVal))
	}
	return values, nil
}

func toHyphenatedString(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func getFlagUsage[T any](flags T, opts options) string {
	f := strings.Builder{}
	w := tabwriter.NewWriter(&f, 0, 0, 4, ' ', 0)
	flagsStruct := reflect.ValueOf(flags)
	for i := 0; i < flagsStruct.NumField(); i++ {
		fd, ok := parseFlagTag(flagsStruct, i)
		if !ok {
			continue
		}
		_, err := fmt.Fprintf(w, "  %s", fd.Flag)
		if err != nil {
			return ""
		}

		if fd.Short != "" {
			_, err := fmt.Fprintf(w, ", %s", fd.Short)
			if err != nil {
				return ""
			}
		}
		if fd.Required {
			_, err := fmt.Fprintf(w, " (required)")
			if err != nil {
				return ""
			}
		}
		if fd.Default != "" {
			_, err := fmt.Fprintf(w, " [default=%s]", fd.Default)
			if err != nil {
				return ""
			}
		}
		if fd.Description != "" {
			_, err := fmt.Fprintf(w, "\t%s", fd.Description)
			if err != nil {
				return ""
			}
		}
		_, err = fmt.Fprintf(w, "\t\n")
		if err != nil {
			return ""
		}
	}

	if opts.HelpFlag != "" {
		_, err := fmt.Fprintf(w, "  %s\tShow help\t\n", opts.HelpFlag)
		if err != nil {
			return ""
		}
	}
	if opts.VersionFlag != "" {
		_, err := fmt.Fprintf(w, "  %s\tShow version\t\n", opts.VersionFlag)
		if err != nil {
			return ""
		}
	}

	err := w.Flush()
	if err != nil {
		return ""
	}
	return f.String()
}

func getSubCommandUsage(cmds []Runnable) string {
	f := strings.Builder{}
	w := tabwriter.NewWriter(&f, 0, 0, 4, ' ', 0)
	for _, cmd := range cmds {
		_, err := fmt.Fprintf(w, "  %s\t%s\t\n", cmd.Name(), cmd.Description())
		if err != nil {
			return ""
		}
	}

	err := w.Flush()
	if err != nil {
		return ""
	}
	return f.String()
}

func defaultUsage[T any](c Command[T], opts options) string {
	b := strings.Builder{}
	b.WriteString(c.Name())
	b.WriteString("\n\t")
	b.WriteString(c.Description())
	if c.Aliases() != nil && len(c.Aliases()) > 0 {
		b.WriteString("\n\n")
		b.WriteString("Aliases: ")
		b.WriteString(strings.Join(c.Aliases(), ", "))
	}
	b.WriteString("\n\nUsage:\n\t")
	b.WriteString(c.Name())
	b.WriteString(" ")
	for _, arg := range c.Args() {
		if arg.Required {
			b.WriteString("<")
			b.WriteString(arg.Name)
			b.WriteString("> ")
		} else {
			b.WriteString("[")
			b.WriteString(arg.Name)
			b.WriteString("] ")
		}
	}
	if reflect.TypeOf(c.Flags()).NumField() > 0 || opts.HelpFlag != "" || opts.VersionFlag != "" {
		b.WriteString("[...options]\n")
	}

	subCommands := c.SubCommands()
	if len(subCommands) > 0 {
		b.WriteString("\t")
		b.WriteString(c.Name())
		b.WriteString(" <command>")
		b.WriteString("\n\nCommands:\n")
		b.WriteString(getSubCommandUsage(subCommands))
	}
	if reflect.TypeOf(c.Flags()).NumField() > 0 || opts.HelpFlag != "" || opts.VersionFlag != "" {
		b.WriteString("\nOptions:\n")
		b.WriteString(getFlagUsage(c.Flags(), opts))
	}
	return b.String()
}
