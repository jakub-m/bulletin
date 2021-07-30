package command

type Command interface {
	Execute(commonOpts Options, args []string) error
}

// Options holds common options that are relevant to all the commands.
type Options struct {
	CacheDir string
}


