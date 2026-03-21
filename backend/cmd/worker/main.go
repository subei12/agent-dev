package main

import (
	"context"
	"log"

	"github.com/your-org/agent-platform/internal/app"
	platformconfig "github.com/your-org/agent-platform/internal/config"
)

func main() {
	cfg := platformconfig.MustLoad()
	worker, err := app.NewWorker(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := worker.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
