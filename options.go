package cligo

import "os"

type options struct {
	Args        []string
	Version     string
	VersionFlag string
	HelpFlag    string
	Strict      bool
	EnvFiles    []string
}

func WithVersion(version string) func(*options) {
	return func(o *options) {
		o.Version = version
	}
}

func WithVersionFlag(versionFlag string) func(*options) {
	return func(o *options) {
		o.VersionFlag = versionFlag
	}
}

func WithHelpFlag(helpFlag string) func(*options) {
	return func(o *options) {
		o.HelpFlag = helpFlag
	}
}

func WithStrict(strict bool) func(*options) {
	return func(o *options) {
		o.Strict = strict
	}
}

func WithEnvFiles(envFiles ...string) func(*options) {
	return func(o *options) {
		o.EnvFiles = envFiles
	}
}

func NewOptions(opts ...func(*options)) options {
	o := options{
		Args:        os.Args[1:],
		Version:     "0",
		VersionFlag: "--version, -v",
		HelpFlag:    "--help, -h",
		Strict:      true,
		EnvFiles:    []string{},
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
