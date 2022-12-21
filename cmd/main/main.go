package main

import (
	"context"
	"pnode/internal/node"
	"pnode/internal/peer"
	"pnode/internal/pool"
)

var ServerNode = &node.Node{
	Addr:  "0.0.0.0:80",
	Peers: make([]*peer.Peer, 0),
	Pool:  pool.New(10),
}

func main() {

	// TODO: graceful shutdown
	ctx := context.Background()
	ServerNode.StartServer(ctx)
}
