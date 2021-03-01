package collector

// Collector is an interface defining a IPFIX collector which collects incoming messages.
type Collector interface {
	Collect()
}

// // ListenAndCollect listens on the UDP network address addr and calls collector to receive incoming IPFIX Messages.
// func ListenAndCollect(addr net.UDPAddr, w io.Writer, bufferSize int, debug bool) error {
// 	c, err := net.ListenUDP(udpNetwork, &addr)
// 	if err != nil {
// 		return errors.Unwrap(err)
// 	}

// 	defer c.Close()

// 	collector := collector{BufferSize: bufferSize, Debug: debug, writer: w}

// 	// decoder := ipfix.NewDecoder(&addr, collector.Debug)
// 	collector.decoder.Source = &addr

// 	for {
// 		buf := make([]byte, collector.BufferSize)
// 		if _, _, err = c.ReadFromUDP(buf); err != nil {
// 			// TODO mateumann: use logger
// 			fmt.Printf("Error occurred: %v\n", err)
// 		}

// 		collector.Collect(buf)
// 	}
// }

// // collector configures IPFIX collection parameters.
// type collector struct {
// 	decoder    *encoding.Decoder
// 	writer     io.Writer // TODO mateumann: replace with ipfix Writer
// 	logger     *log.Logger
// 	BufferSize int
// 	Debug      bool
// }

// func (c *collector) Listen(a *net.UDPAddr) {
// 	// TODO
// }

// // Collect if IPFIX implementation of Collector interface.
// func (c *collector) Collect(p []byte) {
// 	n, err := c.decoder.Read(p)
// 	if err != nil {
// 		c.logger.Printf("Error occurred: %v\n", err)
// 	}

// 	if n > 0 {
// 		// TODO mateumann: replace with ipfix Writer
// 		_, err := c.writer.Write([]byte(c.decoder.String()))
// 		if err != nil {
// 			fmt.Printf("Error occurred: %v\n", err)
// 		}
// 	}
// }

// Collector describes a UDP listener service which listens to the incoming data flow and calls handler function.
// type Collector struct {
// 	address    net.UDPAddr
// 	conn       *net.UDPConn
// 	bufferSize int
// 	debug      bool
// }

// NewUDPCollector starts a UDP server listening on a given IP:port address.
// func NewUDPCollector(ip string, port, bufferSize int, debug bool) (*Collector, error) {
// 	addr := net.UDPAddr{IP: net.ParseIP(ip), Port: port, Zone: ""}

// 	conn, err := net.ListenUDP(collectorNetwork, &addr)
// 	if err != nil {
// 		return nil, errors.Unwrap(err)
// 	}

// 	if debug {
// 		// TODO mateumann: use logger
// 		fmt.Printf("Server listening %v\n", conn.LocalAddr())
// 	}

// 	return &Collector{addr, conn, bufferSize, debug}, nil
// }

// // Close closes conn connection.
// func (c *Collector) Close() error {
// 	return c.conn.Close()
// }

// // Serve registers a io.Writer interface which will store received ipfix messages.
// // It also acts as kind of server function (for time being probably).
// func (c *Collector) Serve(w io.Writer) {
// 	defer c.Close()

// 	decoder := ipfix.NewDecoder(&c.address, c.debug) // TODO mateumann: is `&` needed?

// 	for {
// 		buf := make([]byte, c.bufferSize)
// 		if _, _, err := c.conn.ReadFromUDP(buf); err != nil {
// 			// TODO mateumann: use logger
// 			fmt.Printf("Error occurred: %v\n", err)
// 		}

// 		// var templates []*templateRecord
// 		// var message ipfix.Message

// 		_, err := decoder.Read(buf)
// 		// err = message.UnmarshalBinary(bytes.NewReader(buf), templates)
// 		if err != nil {
// 			log.Print(fmt.Sprintf("%v\n", err))
// 		}

// 		// TODO mateumann: call w.Write somehow on successful read

// 		fmt.Printf("%v\n", decoder)
// 	}
// }

// func (c *Collector) HandleFuncDeprecated(writer io.Writer, reader io.Reader, bufferSize int, debug bool) error {
// 	var udpAddr *net.UDPAddr

// 	for {
// 		buf := make([]byte, bufferSize)
// 		if _, _, err := c.conn.ReadFromUDP(buf); err != nil {
// 			fmt.Printf("Error occurred: %v. Exiting.\n", err)

// 			return errors.Unwrap(err)
// 		}

// 		decoder := ipfix.NewDecoder(udpAddr, debug)

// 		// var templates []*templateRecord
// 		// var message ipfix.Message

// 		_, err := decoder.Read(buf)
// 		// err = message.UnmarshalBinary(bytes.NewReader(buf), templates)
// 		if err != nil {
// 			log.Print(fmt.Sprintf("%v\n", err))
// 		}

// 		fmt.Printf("%v\n", decoder)
// 	}
// }
