package main

import (
	"strings"
)

func filterProfane(msg string) string {
	words := strings.Split(msg, " ")
	for i, word := range words {
		lower := strings.ToLower(word)
		if lower == "kerfuffle" || lower == "sharbert" || lower == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}