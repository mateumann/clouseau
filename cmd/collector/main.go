package main

import (
	"fmt"
	"log"

	"github.com/mateumann/clouseau/net/ipfix"
)

const (
	ipfixIP   = "0.0.0.0"
	ipfixPort = 4739
)

func main() {
	conn, err := ipfix.StartCollector(ipfixIP, ipfixPort)
	if err != nil {
		log.Fatal(err)
	}

	defer ipfix.CloseCollector(conn)

	fmt.Printf("Server listening %v\n", conn.LocalAddr())
}
