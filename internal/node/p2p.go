package node

import (
	"context"
	"fmt"
	"pnode/internal/peer"
)

type Command struct {
	Name  string
	Args  []string
	CType int
	CNode *Node
	CPeer *peer.Peer
}

func (c *Command) Execute(ctx context.Context) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if c.CNode == nil {
		return nil, fmt.Errorf("node is null")
	}

	if c.Name == "" {
		return nil, fmt.Errorf("command name not stated")
	}

	dto := NewCommandDTO(c.CNode, c.CPeer, c.Name, c.Args)
	b, err := p2pCommand(ctx, dto)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Command) Type() int {
	return c.CType
}

func (c *Command) Node() *Node {
	return c.CNode
}

func (c *Command) Peer() *peer.Peer {
	return c.CPeer
}
