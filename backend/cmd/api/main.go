package main

import (
	"log"
	"net/http"

	platformconfig "github.com/your-org/agent-platform/internal/config"
	platformhttp "github.com/your-org/agent-platform/internal/http"
)

func main() {
	cfg := platformconfig.Load()

	log.Printf("api listening on %s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, platformhttp.NewRouter()); err != nil {
		log.Fatal(err)
	}
}
