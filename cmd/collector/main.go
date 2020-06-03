package main

import (
    "fmt"
    "github.com/mateumann/clouseau/encoding/ipfix"
    "log"
    "net"
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

    var message ipfix.Message
    for {
        buf := make([]byte, 8*1024)
        bufferLen, remote, err := conn.ReadFromUDP(buf[:])
        err = ipfix.Unmarshal(buf, &message)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("received %d bytes from %s; header %v\n", bufferLen, remote, message.Header)
    }

}
