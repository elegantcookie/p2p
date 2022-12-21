package pool

import (
	"fmt"
	"net"
	"pnode/internal/peer"
	"sync"
)

type ConnPool struct {
	connections []*peer.Peer
	limit       int

	m sync.Mutex
}

func (c *ConnPool) TryAdd(fn func() (net.Conn, error)) (*peer.Peer, error) {
	c.m.Lock()

	if len(c.connections)+1 > c.limit {
		return nil, fmt.Errorf("pool limit exceed")
	}
	c.m.Unlock()

	conn, err := fn()
	if err != nil {
		return nil, err
	}
	w := peer.Wrap(conn)

	c.m.Lock()
	c.connections = append(c.connections, w)
	c.m.Unlock()

	return w, nil
}

func New(limit int) *ConnPool {
	if limit < 1 {
		limit = 1
	}

	return &ConnPool{
		connections: make([]*peer.Peer, 0),
		limit:       limit,
	}
}
