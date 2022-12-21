package node

import (
	"context"
	"fmt"
	"pnode/internal/apperror"
	"pnode/internal/p2p"
	"pnode/internal/peer"
)

func p2pCommand(ctx context.Context, dto *CommandDTO) ([]byte, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	switch dto.name {
	case "peers":
		switch len(dto.args) {
		case 0:
			peers := []byte(getPeers(dto.n))
			return peers, nil
		case 2:
			switch dto.args[0] {
			case "befriend":
				befriend(dto.n, dto.p)
				return []byte(""), nil
			default:
				return nil, apperror.InvalidCommandBody
			}
		default:
			switch dto.args[0] {
			case "exec":
				if len(dto.args[1:]) > 1 {
					conf := &executeConfig{Target: dto.args[2]}
					execCMD := ""

					for _, arg := range dto.args[2:] {
						execCMD += arg + " "
					}

					//log.Printf(`command="%s", target="%s", command="%s"`, dto.args[0], dto.args[1], execCMD)
					execute(ctx, dto.n, conf, execCMD)
					return []byte("command executed successfully"), nil
				} else {
					return nil, apperror.InvalidCommandBody
				}

			}
		}
	default:
		return nil, apperror.InvalidCommandBody

	}
	return nil, nil
}

func getPeers(n *Node) string {
	return fmt.Sprintf("%v", n.Peers)
}

func befriend(n *Node, p *peer.Peer) {
	n.Peers = append(n.Peers, p)
}

type executeConfig struct {
	Target string
}

func execute(ctx context.Context, n *Node, c *executeConfig, command string) {
	switch c.Target {
	case "all":
		for _, p := range n.Peers {
			select {
			case <-ctx.Done():
				return
			default:
			}

			p2p.MakeRequest(p, []byte(command), false).Do(ctx)
		}
	case "online":
		for _, p := range n.Peers {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if p.Status() == peer.StatusOnline {
				p2p.MakeRequest(p, []byte(command), false).Do(ctx)
			}
		}
	// TODO: peer ID
	default:
		return
	}
}
