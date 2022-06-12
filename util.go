package snuetl

import "strings"

// See https://en.wikipedia.org/wiki/Filename
var surrogates = [][]string{
	{"/", "／"},
	{"\\", "∖"},
	{"?", "？"},
	{"*", "⁎"},
	{":", "∶"},
	{"|", "ǀ"},
	{"\"", "＂"},
}

func sanitizeFileName(name string) string {
	for _, pair := range surrogates {
		name = strings.ReplaceAll(name, pair[0], pair[1])
	}
	return name
}
