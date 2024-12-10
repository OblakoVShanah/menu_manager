package main

import (
	"context"
	"log"
	"menu_manager/internal/app"
)

func main() {
	ctx := context.Background()

	config, err := app.NewConfig("configs/config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	app, err := app.New(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	if err = app.Setup(ctx, config.DB.DSN, config.BarnURL); err != nil {
		log.Fatal(err)
	}

	if err = app.Start(); err != nil {
		log.Fatal(err)
	}
}
