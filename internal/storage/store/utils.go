package store

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/renniemaharaj/news-go/internal/document"
)

func SanitizeFilename(name string) string {
	replacer := strings.NewReplacer(
		"\\", "", "/", "", ":", "", "*", "", "?", "", "\"", "",
		"<", "", ">", "", "|", "", " ", "_",
	)
	name = strings.ToLower(name)
	name = replacer.Replace(name)
	name = strings.TrimSpace(name)
	name = strings.Trim(name, ".-_")
	return name
}

func StringSliceToMap(list []string) map[string]struct{} {
	set := make(map[string]struct{}, len(list))
	for _, v := range list {
		set[strings.ToLower(v)] = struct{}{}
	}
	return set
}

func StringSliceToEmptyStructMap(list []string) map[string]struct{} {
	set := make(map[string]struct{}, len(list))
	for _, v := range list {
		set[v] = struct{}{}
	}
	return set
}

func StringSliceToSanitizedEmptyStructMap(list []string) map[string]struct{} {
	set := make(map[string]struct{}, len(list))
	for _, v := range list {
		set[SanitizeFilename(v)] = struct{}{}
	}
	return set
}

func BytesToReport(bs []byte) (*document.Report, error) {
	var r *document.Report
	err := json.Unmarshal(bs, &r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// Utility function to check wether a file is json object by info
func FileIsObject(info os.DirEntry) bool {
	if !info.IsDir() && strings.HasSuffix(info.Name(), ".json") {
		return true
	}

	return false
}
