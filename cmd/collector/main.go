package main

import (
	"fmt"
	"log"
	"net"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/mateumann/clouseau/encoding/ipfix"
	"github.com/mateumann/clouseau/net/collector"
)

func setupViper() *viper.Viper {
	v := viper.New()

	v.SetDefault("debug", false)
	v.SetDefault("address", "0.0.0.0")
	v.SetDefault("port", 4739)
	v.SetDefault("buffersize", 8192)
	v.SetEnvPrefix("clu_collector")
	v.AutomaticEnv()

	pflag.Bool("debug", false, "Debug mode.")
	pflag.String("address", "0.0.0.0", "IP address of the interface to listen to.")
	pflag.Int("port", 4739, "The UDP port to listen to.")
	pflag.Int("buffersize", 8192, "Buffer size for the incoming messages.")
	pflag.Parse()

	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Printf("%v", err)

		return nil
	}

	return v
}

func main() {
	v := setupViper()

	if v == nil {
		log.Fatal("Unable to read configuration")
	}

	conn, err := collector.StartCollector(v.GetString("address"), v.GetInt("port"))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %v\n", conn.LocalAddr())

	var udpAddr *net.UDPAddr

	for {
		buf := make([]byte, v.GetInt("buffersize"))
		if _, udpAddr, err = conn.ReadFromUDP(buf); err != nil {
			log.Fatal(err) //nolint:gocritic
		}

		decoder := ipfix.NewDecoder(udpAddr, v.GetBool("debug"))

		// var templates []*templateRecord
		// var message ipfix.Message

		_, err = decoder.Read(buf)
		// err = message.UnmarshalBinary(bytes.NewReader(buf), templates)
		if err != nil {
			log.Print(fmt.Sprintf("%v\n", err))
		}

		fmt.Printf("%v\n", decoder)
	}
}
