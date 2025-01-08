package methods

import (
	"fmt"
	"strings"
)

func KeepFirstTenCharacters(s string, length int) string {
	slen := len(s)
	if slen > length {
		slen = length
	}
	str := strings.Replace(s[:slen], "\\x", " ", -1)
	return fmt.Sprintf("%q", str)
}

func StrInSlice(target string, str_array []string) bool {
	for _, element := range str_array {
		if target == element {
			return true
		}
	}
	return false
}
