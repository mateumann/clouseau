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

    for {
        buf := make([]byte, 8*1024)
        bufferLen, remote, err := conn.ReadFromUDP(buf[:])
        var message ipfix.Message
        err = ipfix.Unmarshal(buf, &message)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("%s from %s (%d byte(s), %d set(s))\n", message.Header, remote.IP, bufferLen, len(message.Sets))
        for i, v := range message.Sets {
            fmt.Printf(" #%01d %s\n", i, v)
        }
    }

}
