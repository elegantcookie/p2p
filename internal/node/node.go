package node

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/netip"
	"os"
	"pnode/internal/apperror"
	"pnode/internal/cmd"
	"pnode/internal/cmd/exc"
	"pnode/internal/p2p"
	"pnode/internal/peer"
	"pnode/internal/pool"
	"pnode/internal/prefix"
	"time"
)

var (
	cmdPrefix = []byte("cmd")
	p2pPrefix = []byte("p2p")
)

type Node struct {
	Addr  string // IP address and port of the node
	Peers []*peer.Peer
	Pool  *pool.ConnPool
}

func (node *Node) StartServer(ctx context.Context) {
	ln, err := net.Listen("tcp", node.Addr)
	tls.NewListener(ln, &tls.Config{ClientAuth: tls.RequestClientCert})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	// TODO: rewrite
	err = node.bootstrap()
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		conn, err := node.Pool.TryAdd(func() (net.Conn, error) {
			return ln.Accept()
		})
		if err != nil {
			log.Println(err)
			continue
		}
		log.Printf("%s connected", conn.RemoteAddr())
		go node.handleConnection(ctx, conn)
	}
}

func (node *Node) handleConnection(ctx context.Context, p *peer.Peer) {

	// Read incoming messages from the client
	buf := make([]byte, 1024)

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, err := p.Read(buf)
		if err != nil {
			log.Printf("failed to read incoming message: %s", err)
			return
		}
		log.Printf("incoming message: %v", string(buf[:n]))
		if bytes.Equal(buf[:n], []byte(commandP2PJoin)) {
			log.Println(p.RemoteAddr().String())
			addr, _ := netip.ParseAddrPort(p.RemoteAddr().String())
			p.SetAddress(addr)
			node.Peers = append(node.Peers, p)
			p2p.MakeRequest(p, []byte("peer added"), false)
			//p.SetOutgoing()
			continue
		} else {
			node.handleInput(ctx, p, buf[:n])
		}
	}
}

func (node *Node) handleInput(ctx context.Context, p *peer.Peer, b []byte) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	commandDTO, err := p2p.NewValidator(b).Validate()
	if err != nil {
		p2p.MakeRequest(p, []byte(err.Error()), false).Do(ctx)
		return
	}

	commander, err := node.commander(p, commandDTO.Type, commandDTO.Name, commandDTO.Args)
	if err != nil {
		p2p.MakeRequest(p, []byte(err.Error()), false).Do(ctx)
		return
	}

	bRes, err := commander.Execute(ctx)
	if err != nil {
		p2p.MakeRequest(p, []byte(err.Error()), false).Do(ctx)
		return
	}

	p2p.MakeRequest(p, bRes, false).Do(ctx)
	return

}

func (node *Node) bootstrap() error {
	bootstrapAddr := os.Getenv("BOOTSTRAP_ADDRESS")
	d := &net.Dialer{}
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
	conn, err := d.DialContext(ctx, "tcp", bootstrapAddr)
	if err != nil {
		return err
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Write([]byte(commandP2PJoin))
	if err != nil {
		return err
	}

	buf := make([]byte, 10*1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("failed to bootstrap: %v", err)
	}
	err = json.Unmarshal(buf[:n], &node.Peers)
	if err != nil {
		log.Fatalf("failed to unmarshal bootstrap data: %v", err)
	}

	return nil
}

func (node *Node) commander(p *peer.Peer, cType int, name string, args []string) (cmd.Command, error) {
	switch cType {
	case prefix.TypeCMD:
		return &exc.Command{
			Name:  name,
			Args:  args,
			P:     p,
			CType: cType,
		}, nil
	case prefix.TypeP2P:
		return &Command{
			Name:  name,
			Args:  args,
			CType: cType,
			CNode: node,
		}, nil
	}
	return nil, apperror.CommandNotInterpreted
}
