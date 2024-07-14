package datatype

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func GetFirstNStr(str string, n int) (string, error) {
	rn := []rune(str)
	if len(rn) > n {
		return string(rn[:n]), nil
	}
	return "", fmt.Errorf("invalid string format or size")
}

func IsList(val interface{}) bool {
	valType := reflect.TypeOf(val).String()

	// extracts the first two chars. it expects to capture `[]` from the `[]string`, or `[]int`, or other similar value
	lstChars, err := GetFirstNStr(valType, 2)
	if err == nil && lstChars == "[]" {
		return true
	} else {
		return false
	}
}

// ToSnakeCase casts camel case into snake case
func ToSnakeCase(camelCasestr string) string {
	snakeCaseStr := matchFirstCap.ReplaceAllString(camelCasestr, "${1}_${2}")
	snakeCaseStr = matchAllCap.ReplaceAllString(snakeCaseStr, "${1}_${2}")
	return strings.ToLower(snakeCaseStr)
}
