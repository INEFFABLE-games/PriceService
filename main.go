package main

import (
	"context"
	"github.com/INEFFABLE-games/PriceService/internal/config"
	"github.com/INEFFABLE-games/PriceService/internal/consumer"
	"github.com/INEFFABLE-games/PriceService/internal/service"
	"os"
	"os/signal"
)

func main() {

	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	priceConsumer := consumer.NewPriceConsumer(cfg)
	priceService := service.NewPriceService(priceConsumer)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		priceService.StartStream(ctx)
	}()

	<-c
	cancel()
	os.Exit(1)
}
