package p2p

import (
	"context"
	"pnode/internal/peer"
	"time"
)

type Request struct {
	p        *peer.Peer
	payload  []byte
	response bool
}

func (r *Request) Do(ctx context.Context) (res *Response, err error) {
	res = &Response{}

	select {
	case <-ctx.Done():
		err = ctx.Err()
		return
	default:
	}

	//r.p.Outgoing() <- r.payload

	var deadline time.Time
	if d, ok := ctx.Deadline(); ok {
		deadline = d
	} else {
		deadline = time.Now().Add(5 * time.Second)
	}

	err = r.p.SetWriteDeadline(deadline)
	if err != nil {
		return
	}

	_, err = r.p.Write(r.payload)
	if err != nil {
		return
	}

	if r.response {
		buf := make([]byte, 10*1024)
		n, rErr := r.p.Read(res.b)
		if rErr != nil {
			return
		}

		res.b = buf[:n]
	}
	return

}

func MakeRequest(p *peer.Peer, payload []byte, response bool) *Request {
	return &Request{
		p:        p,
		payload:  payload,
		response: response,
	}
}
