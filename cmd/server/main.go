package main

import (
	"GenealogyTree/internal/app"
	"context"
	"log"
	"os"
	"os/signal"
)

func main() {
	err := realMain()
	if err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := app.Run(ctx); err != nil {
		return err
	}
	return nil
}
