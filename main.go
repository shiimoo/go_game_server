package main

import "github.com/shiimoo/go_game_server/slog"

func main() {
	slog.Default().SetLogPath(func() string {
		return "log.log"
	})
	slog.Log(1, 2, 3, 4, 5)
	slog.Log(10, 2, 3, 4, 5)
}
