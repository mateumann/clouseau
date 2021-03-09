package collector_test

import (
	"log"
	"net"
	"testing"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/mateumann/clouseau/pkg/encoding"
)

//_nolint:funlen.
func TestDecoder_MessagesLen(t *testing.T) {
	tests := []struct {
		name         string
		want         int
		pcapFilename string
	}{
		// {
		// 	"read one data templates",
		// 	1,
		// 	"data_258.pcap",
		// },
		{
			"read two templates",
			2,
			"templates_258_259.pcap",
		},
	}

	addr := net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 4739, Zone: ""}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			decoder := encoding.NewDecoder(&addr, false)
			messagesRead := 0

			handle, err := pcap.OpenOffline("./testdata/pcap/" + tt.pcapFilename)
			if err != nil {
				t.Errorf("error reading pcap file %v: %v", tt.pcapFilename, err)
			}

			// Loop through packets in file
			packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
			for packet := range packetSource.Packets() {
				pcapData := packet.ApplicationLayer()
				n, err := decoder.Read(pcapData.Payload())
				if err != nil {
					log.Fatal(err)
				}
				messagesRead += n
			}

			if messagesRead != tt.want {
				t.Errorf("read %d messages, want: %d", messagesRead, tt.want)
			}

			handle.Close()
		})
	}
}

/*
	decoder := encoding.NewDecoder(&addr, true)

	for {
		buf := make([]byte, c.bufferSize)
		if bytesCount, fromAddr, err := c.conn.ReadFromUDP(buf); err != nil {
			log.Printf("error while reading from UDP occurred: %v\n", err)
		} else {
			log.Printf("%d bytes read from %v\n", bytesCount, fromAddr)
			c.decode(decoder, buf[:bytesCount])
		}
	}
*/
