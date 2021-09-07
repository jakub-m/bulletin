package command

import (
	"bulletin/log"
)

// RunSequence is a "meta" command that runs other commands. In principle it should allow generating default bulletin
// with no user input.
type RunSequence struct {
	Commands []Command
}

func (c *RunSequence) Execute(args []string) error {
	log.Infof("run default sequence")
	for _, c := range c.Commands {
		if err := c.Execute(args); err != nil {
			return err
		}
	}
	return nil
}
