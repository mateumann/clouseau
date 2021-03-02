package collector

// Collector is an interface defining a IPFIX collector which collects incoming messages.
type Collector interface {
	Collect()
}
