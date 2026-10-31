package main

import (
	"fmt"

	"github.com/dworthen/cligo"
)

type myFlags struct {
	Name       string   `cligo:"required,env=TESTING_NAME,short=n,description=name of the entity"`
	Age        uint     `cligo:"required,short=a,description=age of the entity"`
	Files      []string `cligo:"env=TESTING_FILES,short=f,description=list of files"`
	OldeEnough bool     `cligo:"short=o,description=whether the entity is old enough"`
	Score      float64  `cligo:"env=TESTING_SCORE,short=s,description=score of the entity"`
}

func main() {

	coolCmd, _ := cligo.NewCommand(
		cligo.WithName[struct{}]("cool"),
		cligo.WithDescription[struct{}]("Cool command"),
		cligo.WithArgs[struct{}]("<name> <age>"),
		cligo.WithRun(func(args []string, flags struct{}) error {
			fmt.Printf("Running command with args: %v and flags: %+v\n", args, flags)
			return nil
		}),
	)

	cmd, err := cligo.NewCommand(
		cligo.WithName[myFlags]("start"),
		cligo.WithAliases[myFlags]([]string{"run", "init"}),
		cligo.WithDescription[myFlags]("Start the entity"),
		cligo.WithArgs[myFlags]("<name> [age]"),
		cligo.WithFlags(myFlags{
			Age:   30,
			Files: []string{"file1.txt", "file2.txt"},
		}),
		cligo.WithRun(func(args []string, flags myFlags) error {
			fmt.Printf("Running command with args: %v and flags: %+v\n", args, flags)
			return nil
		}),
		cligo.WithSubCommands[myFlags](coolCmd),
	)
	if err != nil {
		fmt.Printf("error creating command: %s\n", err)
		return
	}

	cligo.MustRunCommandWithOptions(cmd, cligo.NewOptions(cligo.WithStrict(true), cligo.WithEnvFiles(".env")))
}
