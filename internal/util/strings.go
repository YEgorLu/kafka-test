package util

import "strings"

func PadPrefix(s string, length int, symb string) string {
	if len(s) >= length {
		return s
	}
	padString := strings.Repeat(symb, length-len(s))
	return padString + s
}
