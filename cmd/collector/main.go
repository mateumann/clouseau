package main

import (
	"fmt"
	"log"
	"net"

	"github.com/mateumann/clouseau/encoding/ipfix"
)

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 4739,
		IP:   net.ParseIP("0.0.0.0"),
	})

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %s\n", conn.LocalAddr().String())

	var templates *[]ipfix.TemplateRecord

	for {
		buf := make([]byte, 8*1024)
		bufferLen, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Fatal(err)
		}
		var message ipfix.Message
		if err := ipfix.Unmarshal(buf, templates, &message); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s from %s (%d byte(s), %d set(s))\n", message.Header, remote.IP, bufferLen, len(message.Sets))
		for i, v := range message.Sets {
			fmt.Printf(" #%01d %s\n", i, v)
		}
	}
}
