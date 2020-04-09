package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 2055,
		IP:   net.ParseIP("0.0.0.0"),
	})

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %s\n", conn.LocalAddr().String())

	for {
		message := make([]byte, 1024)
		rlen, remote, err := conn.ReadFromUDP(message[:])
		if err != nil {
			log.Fatal(err)
		}
		// // Handle the connection in a new goroutine.
		// // The loop then returns to accepting, so that
		// // multiple connections may be served concurrently.
		// go func(c net.Conn) {
		// 	// Echo all incoming data.
		// 	io.Copy(c, c)
		// 	// Shut down the connection.
		// 	c.Close()
		// }(conn)
		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Printf("received %d bytes (%s) from %s\n", rlen, data, remote)
	}

}
