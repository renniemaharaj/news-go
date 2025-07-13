package coordinator

import "strings"

func sanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"\\", "", "/", "", ":", "", "*", "", "?", "", "\"", "",
		"<", "", ">", "", "|", "", " ", "_",
	)
	return replacer.Replace(name)
}
