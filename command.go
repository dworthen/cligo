package cligo

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type argDescriptor struct {
	Name     string
	Required bool
}

type command[T any] struct {
	name        string
	aliases     []string
	description string
	flags       T
	args        []argDescriptor
	run         func(args []string, flags T) error
	usage       func(cmd Command[T], opts options) string
	subCommands []Runnable
}

func (c *command[T]) SubCommands() []Runnable {
	return c.subCommands
}

func (c *command[T]) Name() string {
	return c.name
}

func (c *command[T]) Description() string {
	return c.description
}

func (c *command[T]) Aliases() []string {
	return c.aliases
}

func (c *command[T]) Flags() T {
	return c.flags
}

func (c *command[T]) Args() []argDescriptor {
	return c.args
}

func (c *command[T]) Usage(opts options) string {
	if c.usage == nil {
		return "Usage not implemented"
	}
	return c.usage(c, opts)
}

func (c *command[T]) Execute(args []string, opts options) error {

	cmdArgs := gatherArgs(args)
	argDescriptors := c.Args()
	if opts.Strict && len(cmdArgs) > len(argDescriptors) {
		return fmt.Errorf("too many arguments provided %v", cmdArgs)
	}
	if len(cmdArgs) < len(argDescriptors) {
		nextArgDescriptor := argDescriptors[len(cmdArgs)]
		if nextArgDescriptor.Required {
			return fmt.Errorf("missing required argument: %s", nextArgDescriptor.Name)
		}
	}
	flagValues := args[len(cmdArgs):]

	flags, err := c.getFlags(flagValues, opts)
	if err != nil {
		return err
	}
	return c.Run(cmdArgs, flags)
}

func (c *command[T]) Run(args []string, flags T) error {
	if c.run == nil {
		return fmt.Errorf("command not implemented")
	}
	return c.run(args, flags)
}

func (c *command[T]) getFlags(flagArgs []string, opts options) (T, error) {
	flags := c.Flags()

	allFlags := gatherFlags(flagArgs)

	flagsStruct := reflect.ValueOf(&flags).Elem()
	for i := 0; i < flagsStruct.NumField(); i++ {
		field := flagsStruct.Type().Field(i)
		value := flagsStruct.Field(i)
		strType := field.Type.String()
		fd, ok := parseFlagTag(flagsStruct, i)
		if !ok {
			continue
		}

		strSubType := strType
		isSlice := strings.HasPrefix(strType, "[]")
		if isSlice {
			strSubType = strType[2:]
		}

		switch strSubType {
		case "string":
			flagValues, err := getStringFlags(flagArgs, []string{fd.Flag, fd.Short})
			if err != nil {
				return flags, err
			}

			if len(flagValues) == 0 && fd.Env != "" {
				envValue := os.Getenv(fd.Env)
				if envValue != "" {
					flagValues = strings.Split(envValue, ",")
				}
			}

			if len(flagValues) > 0 {
				allFlags[fd.Flag] = true
				allFlags[fd.Short] = true
				if isSlice {
					flagsStruct.FieldByName(field.Name).Set(reflect.ValueOf(flagValues))
				} else {
					flagsStruct.FieldByName(field.Name).SetString(flagValues[0])
				}
			}
		case "int":
			strFlagValues, err := getStringFlags(flagArgs, []string{fd.Flag, fd.Short})
			if err != nil {
				return flags, err
			}

			if len(strFlagValues) == 0 && fd.Env != "" {
				envValue := os.Getenv(fd.Env)
				if envValue != "" {
					strFlagValues = strings.Split(envValue, ",")
				}
			}

			flagValues, err := getIntFlags(strFlagValues, fd.Flag)
			if err != nil {
				return flags, err
			}

			if len(flagValues) > 0 {
				allFlags[fd.Flag] = true
				allFlags[fd.Short] = true
				if isSlice {
					flagsStruct.FieldByName(field.Name).Set(reflect.ValueOf(flagValues))
				} else {
					flagsStruct.FieldByName(field.Name).SetInt(int64(flagValues[0]))
				}
			}
		case "uint":
			strFlagValues, err := getStringFlags(flagArgs, []string{fd.Flag, fd.Short})
			if err != nil {
				return flags, err
			}

			if len(strFlagValues) == 0 && fd.Env != "" {
				envValue := os.Getenv(fd.Env)
				if envValue != "" {
					strFlagValues = strings.Split(envValue, ",")
				}
			}

			flagValues, err := getUintFlags(strFlagValues, fd.Flag)
			if err != nil {
				return flags, err
			}

			if len(flagValues) > 0 {
				allFlags[fd.Flag] = true
				allFlags[fd.Short] = true
				if isSlice {
					flagsStruct.FieldByName(field.Name).Set(reflect.ValueOf(flagValues))
				} else {
					flagsStruct.FieldByName(field.Name).SetUint(uint64(flagValues[0]))
				}
			}
		case "float32":
			strFlagValues, err := getStringFlags(flagArgs, []string{fd.Flag, fd.Short})
			if err != nil {
				return flags, err
			}

			if len(strFlagValues) == 0 && fd.Env != "" {
				envValue := os.Getenv(fd.Env)
				if envValue != "" {
					strFlagValues = strings.Split(envValue, ",")
				}
			}

			flagValues, err := getFloat32Flags(strFlagValues, fd.Flag)
			if err != nil {
				return flags, err
			}

			if len(flagValues) > 0 {
				allFlags[fd.Flag] = true
				allFlags[fd.Short] = true
				if isSlice {
					flagsStruct.FieldByName(field.Name).Set(reflect.ValueOf(flagValues))
				} else {
					flagsStruct.FieldByName(field.Name).SetFloat(float64(flagValues[0]))
				}
			}
		case "float64":
			strFlagValues, err := getStringFlags(flagArgs, []string{fd.Flag, fd.Short})
			if err != nil {
				return flags, err
			}

			if len(strFlagValues) == 0 && fd.Env != "" {
				envValue := os.Getenv(fd.Env)
				if envValue != "" {
					strFlagValues = strings.Split(envValue, ",")
				}
			}

			flagValues, err := getFloat64Flags(strFlagValues, fd.Flag)
			if err != nil {
				return flags, err
			}

			if len(flagValues) > 0 {
				allFlags[fd.Flag] = true
				allFlags[fd.Short] = true
				if isSlice {
					flagsStruct.FieldByName(field.Name).Set(reflect.ValueOf(flagValues))
				} else {
					flagsStruct.FieldByName(field.Name).SetFloat(float64(flagValues[0]))
				}
			}
		case "bool":
			flagValues, err := getBoolFlags(flagArgs, []string{fd.Flag, fd.Short})
			if err != nil {
				return flags, err
			}

			if len(flagValues) == 0 && fd.Env != "" {
				envValue := os.Getenv(fd.Env)
				if envValue != "" {
					parts := strings.Split(envValue, ",")
					for _, part := range parts {
						flagValues = append(flagValues, strings.ToLower(part) == "true")
					}
				}
			}

			if len(flagValues) > 0 {
				allFlags[fd.Flag] = true
				allFlags[fd.Short] = true
				if isSlice {
					flagsStruct.FieldByName(field.Name).Set(reflect.ValueOf(flagValues))
				} else {
					flagsStruct.FieldByName(field.Name).SetBool(flagValues[0])
				}
			}
		default:
			return flags, fmt.Errorf("unsupported flag type: %s", strType)
		}

		if fd.Required && strSubType != "bool" {
			if value.IsZero() {
				return flags, fmt.Errorf("missing required flag: %s", fd.Flag)
			}
		}
	}

	if opts.Strict {
		for k, v := range allFlags {
			if !v {
				return flags, fmt.Errorf("unknown flag: %s", k)
			}
		}
	}

	return flags, nil
}

func WithName[T any](name string) func(*command[T]) error {
	return func(c *command[T]) error {
		c.name = name
		return nil
	}
}

func WithAliases[T any](aliases []string) func(*command[T]) error {
	return func(c *command[T]) error {
		c.aliases = aliases
		return nil
	}
}

func WithDescription[T any](description string) func(*command[T]) error {
	return func(c *command[T]) error {
		c.description = description
		return nil
	}
}

func WithFlags[T any](flags T) func(*command[T]) error {

	return func(c *command[T]) error {
		if reflect.TypeOf(flags).Kind() != reflect.Struct {
			return fmt.Errorf("flags must be a pointer to a struct")
		}
		c.flags = flags
		return nil
	}
}

func WithArgs[T any](args string) func(*command[T]) error {
	return func(c *command[T]) error {
		argsSlice := strings.Split(args, " ")

		argDescriptors := make([]argDescriptor, len(argsSlice))
		var hasOptional bool
		for i, arg := range argsSlice {
			hasOptional = hasOptional || strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]")
			isRequired := strings.HasPrefix(arg, "<") && strings.HasSuffix(arg, ">")
			if hasOptional && isRequired {
				return fmt.Errorf("required argument cannot follow an optional argument")
			}
			// strip angle and square brackets from the argument name
			arg = strings.TrimPrefix(arg, "<")
			arg = strings.TrimSuffix(arg, ">")
			arg = strings.TrimPrefix(arg, "[")
			arg = strings.TrimSuffix(arg, "]")

			argDescriptors[i] = argDescriptor{
				Name:     arg,
				Required: isRequired,
			}
		}

		c.args = argDescriptors
		return nil
	}
}

func WithSubCommands[T any](subCommands ...Runnable) func(*command[T]) error {
	return func(c *command[T]) error {
		c.subCommands = subCommands
		return nil
	}
}

func WithRun[T any](run func(args []string, flags T) error) func(*command[T]) error {
	return func(c *command[T]) error {
		c.run = run
		return nil
	}
}

func NewCommand[T any](opts ...func(*command[T]) error) (Command[T], error) {

	cmd := &command[T]{
		usage: defaultUsage[T],
	}

	for _, opt := range opts {
		if err := opt(cmd); err != nil {
			return nil, err
		}
	}

	if strings.TrimSpace(cmd.Name()) == "" {
		return nil, fmt.Errorf("command name is required. Specify command name using WithName")
	}

	return cmd, nil
}

func MustNewCommand[T any](opts ...func(*command[T]) error) Command[T] {
	cmd, err := NewCommand(opts...)
	if err != nil {
		panic(err)
	}
	return cmd
}
