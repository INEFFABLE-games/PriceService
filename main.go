package main

import (
	"context"
	"fmt"
	"github.com/INEFFABLE-games/PriceService/internal/config"
	"github.com/INEFFABLE-games/PriceService/internal/consumer"
	"github.com/INEFFABLE-games/PriceService/internal/server"
	"github.com/INEFFABLE-games/PriceService/internal/service"
	"github.com/INEFFABLE-games/PriceService/protocol"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
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

	channels := map[int]chan []byte{}

	// Stats new grpc server
	go func() {
		grpcServer := grpc.NewServer()
		pricesServer := server.NewPriceServer(ctx, channels)
		protocol.RegisterPriceServiceServer(grpcServer, pricesServer)

		lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%v", cfg.GrpcPort))
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		err = grpcServer.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Starts new redis stream
	go func() {
		priceService.StartStream(ctx, channels)
	}()

	<-c
	cancel()
	os.Exit(1)
}
