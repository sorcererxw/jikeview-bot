package util

import "regexp"

var spaceAndPunctuationRegexp = regexp.MustCompile("[\\s,.|~!?@#$%/^&*_+=\"'“”\\[\\]{}()<>「」《》。，？！]")

func RemoveAllSpaceAndPunctuation(s string) string {
	return spaceAndPunctuationRegexp.ReplaceAllString(s, "")
}
