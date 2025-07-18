package utils

import "strings"

func EmptyMapToStringSlice(m map[string]struct{}) []string {
	list := make([]string, 0, len(m))
	for k := range m {
		list = append(list, k)
	}
	return list
}

func StringSliceToMap(list []string) map[string]struct{} {
	set := make(map[string]struct{}, len(list))
	for _, v := range list {
		set[strings.ToLower(v)] = struct{}{}
	}
	return set
}
