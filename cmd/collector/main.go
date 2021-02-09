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

type config struct {
	debug         bool   `yaml:"debug"`
	listenAddress string `yaml:"address"`
	listenPort    int    `yaml:"port"`
	bufferSize    int    `yaml:"buffer-size"`
}

func parseFlags() *config {
	var c config

	pflag.Bool("debug", false, "Debug mode.")
	pflag.String("address", "0.0.0.0", "IP address of the interface to listen to.")
	pflag.Int("port", 4739, "The UDP port to listen to.")
	pflag.Int("buffer-size", 8192, "Buffer size for the incoming messages.")
	pflag.Parse()

	v := viper.New()

	if err := v.BindPFlags(pflag.CommandLine); err != nil {
		fmt.Printf("%v", err)

		return nil
	}

	v.AddConfigPath(".")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		println("ReadInConfig error %v.", err)

		return nil
	}

	if err := v.Unmarshal(&c); err != nil {
		println("Unmarshal error %v.", err)

		return nil
	}

	return &c
}

func main() {
	args := parseFlags()

	if args == nil {
		log.Fatal("Unable to read configuration")
	}

	conn, err := collector.StartCollector(args.listenAddress, args.listenPort)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("Server listening %v\n", conn.LocalAddr())

	var udpAddr *net.UDPAddr

	for {
		buf := make([]byte, args.bufferSize)
		if _, udpAddr, err = conn.ReadFromUDP(buf); err != nil {
			log.Fatal(err) //nolint:gocritic
		}

		decoder := ipfix.NewDecoder(udpAddr, args.debug)

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
