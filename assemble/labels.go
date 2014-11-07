package assemble

import (
	"regexp"
)

func tryLabel(line string) string {
	re := regexp.MustCompile("([A-Za-z]+):")
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return ""
	}
	return matches[1]
}
