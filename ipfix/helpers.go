package ipfix

import "fmt"

func intDescribed(n int, singular, plural string) string {
	if n == 1 {
		return fmt.Sprintf("1 %v", singular)
	}

	return fmt.Sprintf("%v %v", n, plural)
}

func fieldDescribed(n uint16) string {
	return intDescribed(int(n), "field", "fields")
}
