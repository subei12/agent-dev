package main

import (
	"context"
	"log"

	"github.com/your-org/agent-platform/internal/app"
	platformconfig "github.com/your-org/agent-platform/internal/config"
)

func main() {
	cfg := platformconfig.MustLoad()
	api, err := app.NewAPI(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("api listening on %s", cfg.Addr)
	if err := api.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
