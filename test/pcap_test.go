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
		pcapFilename string
		want         int
	}{
		{
			"read two templates",
			"templates_258_259.pcap",
			1,
		},
		{
			"read one data record",
			"data_258.pcap",
			1,
		},
		{
			"read six data record",
			"six_data_258.pcap",
			6,
		},
		{
			"templates and data ",
			"templates_and_data_258_259.pcap",
			3,
		},
	}

	addr := net.UDPAddr{IP: net.ParseIP(""), Port: 0, Zone: ""}

	for _, tt := range tests {
		tt := tt // pin

		t.Run(tt.name, func(t *testing.T) {
			decoder := encoding.NewDecoder(&addr, true)
			messagesRead := 0

			handle, err := pcap.OpenOffline("./testdata/pcap/" + tt.pcapFilename)
			if err != nil {
				t.Errorf("error reading pcap file %v: %v", tt.pcapFilename, err)
			}

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
