package util

import "strings"

func CapitalizeFirst(str string) string {
	if str == "" {
		return ""
	}
	return strings.ToUpper(str[:1]) + str[1:]
}
