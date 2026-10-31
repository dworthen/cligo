package cligo

type Runnable interface {
	Name() string
	Aliases() []string
	Description() string
	Args() []argDescriptor
	Usage(opts options) string
	SubCommands() []Runnable
	Execute(args []string, opts options) error
}

type Command[T any] interface {
	Runnable
	Flags() T
	Run(args []string, flags T) error
}
