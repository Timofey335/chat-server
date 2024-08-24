package main

import (
	"context"
	"flag"
	"log"

	"github.com/Timofey335/chat-server/internal/app"
	"github.com/Timofey335/chat-server/internal/config"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	if err := config.Load(configPath); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	a, err := app.NewApp(ctx, configPath)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	if err := a.Run(); err != nil {
		log.Fatalf("failed ti run app: %s", err.Error())
	}
}
