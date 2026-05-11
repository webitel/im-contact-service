package main

import (
	"log/slog"

	"github.com/webitel/im-contact-service/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		slog.Error("[MAIN] running server", "error", err)

		return
	}
}
