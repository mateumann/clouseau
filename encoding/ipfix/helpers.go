package ipfix

import "fmt"

func intDescribed(n int, singular string, plural string) string {
	if n == 1 {
		return fmt.Sprintf("1 %v", singular)
	}

	return fmt.Sprintf("%v %v", n, plural)
}

func uint16Described(n uint16, singular string, plural string) string {
	return intDescribed(int(n), singular, plural)
}
