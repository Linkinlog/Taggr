package main

import (
	"log/slog"
)

func main() {
	g := NewGame()

	if err := ServeHTTP(g); err != nil {
		slog.Error(err.Error())
	}
}
