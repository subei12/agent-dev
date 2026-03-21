package main

import (
	"context"
	"log"

	"github.com/your-org/agent-platform/internal/app"
	platformconfig "github.com/your-org/agent-platform/internal/config"
)

func main() {
	cfg := platformconfig.MustLoad()
	worker := app.NewWorker(cfg)

	if err := worker.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
