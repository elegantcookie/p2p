package node

import "pnode/internal/peer"

type CommandDTO struct {
	n    *Node
	p    *peer.Peer
	name string
	args []string
}

func NewCommandDTO(n *Node, p *peer.Peer, name string, args []string) *CommandDTO {
	return &CommandDTO{
		n:    n,
		p:    p,
		name: name,
		args: args,
	}
}
