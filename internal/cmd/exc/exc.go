package exc

import (
	"context"
	"os/exec"
	"pnode/internal/apperror"
	"pnode/internal/peer"
)

type Command struct {
	Name string
	Args []string

	CType int
	P     *peer.Peer
}

func (c *Command) Execute(ctx context.Context) ([]byte, error) {
	if c.Name == "" {
		return nil, apperror.InvalidCommandBody
	}
	if c.Args == nil || len(c.Args) == 0 {
		return exec.Command(c.Name).Output()
	}
	return exec.CommandContext(ctx, c.Name, c.Args...).Output()
}

func (c *Command) Type() int {
	return c.CType
}
