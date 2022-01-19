package main

import "fmt"

func format_json(t OutputResult) string {
	return `{"target": "` + t.target + `", "version": "` + t.version + `", "players": "` + t.players + `", "description": "` + t.description + `"}`
}

func format_csv(t OutputResult) string {
	return `"` + t.target + `","` + t.version + `","` + t.players + `","` + t.description + `"`
}

func format_qubo(t OutputResult) string {
	return fmt.Sprintf("(%s)(%s)(%s)(%s)", t.target, t.players, t.version, t.description)
}
