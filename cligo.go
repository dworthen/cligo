package cligo

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/dworthen/cligo/internal/util"
)

var once sync.Once
var envErr error

func RunCommand(cmd Runnable, options ...func(*options)) (Runnable, error) {
	opts := NewOptions(options...)

	return runCommandWithOptions(cmd, opts)
}

func runCommandWithOptions(cmd Runnable, opts options) (Runnable, error) {

	once.Do(func() {
		for _, envFile := range opts.EnvFiles {
			if err := util.LoadEnv(envFile); err != nil {
				envErr = err
				break
			}
		}
	})

	if envErr != nil {
		return cmd, envErr
	}

	if opts.VersionFlag != "" {
		vFlags := strings.Split(opts.VersionFlag, ",")
		for _, vFlag := range vFlags {
			if slices.Contains(opts.Args, strings.TrimSpace(vFlag)) {
				fmt.Println(opts.Version)
				return nil, nil
			}
		}
	}

	for index, arg := range opts.Args {
		if strings.HasPrefix(arg, "-") {
			break
		}
		subCommands := cmd.SubCommands()
		for _, subCmd := range subCommands {
			if subCmd.Name() == arg || slices.Contains(subCmd.Aliases(), arg) {
				opts.Args = opts.Args[index+1:]
				return runCommandWithOptions(subCmd, opts)
			}
		}
	}

	if opts.HelpFlag != "" {
		hFlags := strings.Split(opts.HelpFlag, ",")
		for _, hFlag := range hFlags {
			if slices.Contains(opts.Args, strings.TrimSpace(hFlag)) {
				fmt.Println(cmd.Usage(opts))
				return nil, nil
			}
		}
	}
	return cmd, cmd.Execute(opts.Args, opts)
}

func MustRunCommand(cmd Runnable, options ...func(*options)) {

	opts := NewOptions(options...)

	actualCmd, err := runCommandWithOptions(cmd, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error running %s: %s\n", actualCmd.Name(), err)
		fmt.Fprintf(os.Stderr, "%s\n", actualCmd.Usage(opts))
		os.Exit(1)
	}
}
