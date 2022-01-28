package utils

import (
	"fmt"
)

func FormatJSON(t OutputResult) string {
	return `{"target": "` + t.Target + `", "version": "` + t.Version + `", "players": "` + t.Players + `", "description": "` + t.Description + `"}`
}

func FormatCSV(t OutputResult) string {
	return `"` + t.Target + `","` + t.Version + `","` + t.Players + `","` + t.Description + `"`
}

func FormatQubo(t OutputResult) string {
	return fmt.Sprintf("(%s)(%s)(%s)(%s)", t.Target, t.Players, t.Version, t.Description)
}
