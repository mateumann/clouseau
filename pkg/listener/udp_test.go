package listener

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"syscall"
	"testing"
)

// nolint: funlen
func TestNewUDPListener(t *testing.T) {
	testCases := []struct {
		addr string
		want *net.UDPAddr
		err  error
	}{
		{
			addr: "127.0.0.1:123",
			want: &net.UDPAddr{
				IP:   net.IP{127, 0, 0, 1},
				Port: 123,
			},
			err: nil,
		},
		{
			addr: "[fe80:1234:5678::9]:12345",
			want: &net.UDPAddr{
				IP:   net.IP{0xfe, 0x80, 0x12, 0x34, 0x56, 0x78, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x09},
				Port: 12345,
			},
			err: nil,
		},
		{
			addr: "0.0.0.0:9876",
			want: &net.UDPAddr{
				IP:   net.IPv4zero,
				Port: 9876,
			},
			err: nil,
		},
		{
			addr: "[::0]:9876",
			want: &net.UDPAddr{
				IP:   net.IPv6zero,
				Port: 9876,
			},
			err: nil,
		},
		{
			addr: "[::0]:abcdefgh",
			want: nil,
			err: &Error{
				op: "create new UDPListener",
				err: &strconv.NumError{
					Func: "Atoi",
					Num:  "abcdefgh",
					Err:  errors.New("invalid syntax"), // nolint: goerr113
				},
			},
		},
		{
			addr: "ijkl",
			want: nil,
			err: &Error{
				op: "create new UDPListener",
				err: &net.AddrError{
					Addr: "ijkl",
					Err:  "missing port in address",
				},
			},
		},
	}

	for _, tc := range testCases {
		l, err := NewUDPListener(tc.addr)

		if err != nil || tc.err != nil {
			if fmt.Sprint(err) != fmt.Sprint(tc.err) {
				t.Errorf("unexpected error occurred NewUDPListener(%s): %v; expected %v", tc.addr, err, tc.err)
			}

			if fmt.Sprint(errors.Unwrap(err)) != fmt.Sprint(errors.Unwrap(tc.err)) {
				t.Errorf("unexpected wrapped error occurred NewUDPListener(%s): %v; expected %v", tc.addr, err, tc.err)
			}
		} else if got := *l.listenAddress; !got.IP.Equal(tc.want.IP) || got.Port != tc.want.Port ||
			got.Zone != tc.want.Zone {
			t.Errorf("NewUDPListener(%s) = %v; want %v", tc.addr, got, tc.want)
		}
	}
}

func TestUDPListener_Listen(t *testing.T) {
	testCases := []struct {
		num int
		err error
	}{
		{
			num: -100,
			err: &Error{
				op:  "invalid number of concurrent listeners: -100",
				err: nil,
			},
		},
		{
			num: 0,
			err: &Error{
				op:  "invalid number of concurrent listeners: 0",
				err: nil,
			},
		},
		{
			num: 1,
			err: &Error{
				op: "listen UDP",
				err: &net.OpError{
					Op:     "listen",
					Net:    "udp",
					Source: nil,
					Addr: &net.UDPAddr{
						IP:   net.IPv4(127, 0, 0, 1),
						Port: 12,
						Zone: "",
					},
					Err: &os.SyscallError{
						Syscall: "bind",
						Err:     syscall.EACCES,
					},
				},
			},
		},
	}

	l, err := NewUDPListener("127.0.0.1:12")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range testCases {
		err = l.Listen(tc.num)

		if err != nil || tc.err != nil {
			if fmt.Sprint(err) != fmt.Sprint(tc.err) {
				t.Errorf("unexpected error occurred UDPListener.Listen(%d): wanted %v, got %v", tc.num, tc.err, err)
			}
		}
	}
}
