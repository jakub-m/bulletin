package command

type Command interface {
	Execute(args []string) error
}
