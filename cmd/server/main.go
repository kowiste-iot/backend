package main

import (
	"backend/internal/core"
	"context"
)

func main() {
	ctx := context.Background()
	app, err := core.NewCore(ctx)
	if err != nil {
		panic(err)
	}

	if err := app.Start(ctx); err != nil {
		panic(err)
	}
}
