package travelskypnr

import "regexp"

var reg = regexp.MustCompile(`\d{2}.RMK\s?\w{2}\/(\w{6})`)

func GetBigCode(pnr string) string {

	matches := reg.FindStringSubmatch(pnr)
	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}
