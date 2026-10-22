package main

import (
	"log/slog"

	"github.com/webitel/webitel-go-kit/pkg/semconv"

	"github.com/webitel/im-contact-service/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		slog.Error("[MAIN] running server", semconv.ErrorKey, err)

		return
	}
}
