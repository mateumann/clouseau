package listener

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"time"
)

type udpHandlerCtxKey uint8

const (
	// UDPSourceAddress is the key which describes a UDP source address value in the context.Context map.
	UDPSourceAddress udpHandlerCtxKey = iota // UDPSourceAddress = 0
)

type UDPListener struct {
	listenAddress      *net.UDPAddr
	listenTimeout      time.Duration
	expectedSrcAddress *net.UDPAddr
	bufferSize         int
	handlers           []Handler
}

func (l *UDPListener) Listen(concurrentListeners int) error {
	if concurrentListeners <= 0 {
		return &Error{
			op:  fmt.Sprintf("invalid number of concurrent listeners: %d", concurrentListeners),
			err: nil,
		}
	}

	conn, err := net.ListenUDP(l.listenAddress.Network(), l.listenAddress)
	if err != nil {
		return &Error{op: "listen UDP", err: err}
	}

	defer conn.Close()

	if l.listenTimeout != 0 {
		if err = conn.SetReadDeadline(time.Now().Add(l.listenTimeout)); err != nil {
			return &Error{op: "set read deadline", err: err}
		}
	}

	quit := make(chan struct{})

	for i := 0; i < concurrentListeners; i++ {
		go l.listen(conn, quit)
	}
	<-quit

	return nil
}

func (l *UDPListener) listen(conn *net.UDPConn, quit chan struct{}) {
	buf := make([]byte, l.bufferSize)

	var ctx *context.Context

	n, err := 0, error(nil)

	for err == nil {
		n, ctx, err = l.read(conn, buf)

		if err != nil || ctx == nil {
			continue
		}

		for _, h := range l.handlers {
			h.Setup(ctx)

			go h.Handle(buf[:n])
		}
	}
	quit <- struct{}{}
}

func (l *UDPListener) read(conn *net.UDPConn, b []byte) (int, *context.Context, error) {
	n, remoteAddress, err := conn.ReadFromUDP(b)
	if err != nil {
		return n, nil, &Error{op: "UDPListener read", err: err}
	}

	if l.expectedSrcAddress != nil && !UDPAddressEqual(*l.expectedSrcAddress, *remoteAddress) {
		return n, nil, nil
	}

	ctx := context.WithValue(context.Background(), UDPSourceAddress, remoteAddress)

	return n, &ctx, nil
}

func NewUDPListener(listenAddress string) (*UDPListener, error) {
	host, port, err := net.SplitHostPort(listenAddress)
	if err != nil {
		return nil, &Error{op: "create new UDPListener", err: err}
	}

	var portNum int
	portNum, err = strconv.Atoi(port)

	if err != nil {
		return nil, &Error{op: "create new UDPListener", err: err}
	}

	return &UDPListener{
		listenAddress: &net.UDPAddr{
			IP:   net.ParseIP(host),
			Port: portNum,
			Zone: "",
		},
		listenTimeout:      0,
		expectedSrcAddress: nil,
		bufferSize:         0,
		handlers:           nil,
	}, nil
}
