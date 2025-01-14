package main

import (
	"context"
	"log"
	"ddd/internal/core"
)

func main() {
	ctx := context.Background()
	
	app, err := core.NewCore(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	
	if err := app.Start(ctx); err != nil {
		log.Fatalf("Failed to start application: %v", err)
	}
}