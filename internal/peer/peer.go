package peer

import (
	"fmt"
	"net"
	"net/netip"
)

type Peer struct {
	net.Conn
	status  int
	address netip.AddrPort

	incoming []byte
	outgoing []byte
}

func resolveStatus(status int) string {
	switch status {
	case StatusOnline:
		return "online"
	case StatusOffline:
		return "offline"
	default:
		return "unauthorized"
	}
}

func (p *Peer) Close() error {
	p.status = StatusOffline
	return p.Conn.Close()
}

func (p *Peer) Outgoing() []byte {
	return p.outgoing
}

func (p *Peer) Incoming() []byte {
	return p.incoming
}

func (p *Peer) SetOutgoing(b []byte) {
	p.outgoing = b

}

func (p *Peer) SetIncoming(b []byte) {
	p.incoming = b
}

func (p *Peer) Status() int {
	return p.status
}

func (p Peer) ClearIncoming() {
	p.incoming = make([]byte, 0)
}

func (p Peer) ClearOutgoing() {
	p.outgoing = make([]byte, 0)
}

func (p Peer) Address() netip.AddrPort {
	return p.address
}

func (p Peer) SetAddress(a netip.AddrPort) {
	p.address = a
}

func (p Peer) String() string {
	return fmt.Sprintf(`{ address="%s", status="%s" }`, p.address, resolveStatus(p.status))
}

func Wrap(conn net.Conn) *Peer {
	return &Peer{
		Conn:     conn,
		status:   StatusOnline,
		incoming: make([]byte, 0),
		outgoing: make([]byte, 0),
	}
}
