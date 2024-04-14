package checkdigits

import (
	"regexp"
)

func CheckInSring(s string) (bool, error) {
	pattern := `^\d*$`
	matched, err := regexp.Match(pattern, []byte(s))
	if err != nil || !matched {
		return false, err
	}
	return true, nil
}
