package service

import (
	"context"
	"encoding/json"
	"github.com/INEFFABLE-games/PriceService/internal/consumer"
	"github.com/INEFFABLE-games/PriceService/internal/protocol"
	"github.com/INEFFABLE-games/PriceService/internal/server"
	"github.com/INEFFABLE-games/PriceService/models"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"time"
)

// PriceService is structure for PriceService object
type PriceService struct {
	consumer *consumer.PriceConsumer
	server   *protocol.PriceServiceServer
}

// StartStream starts infinity cycle to get butch of prices from redis
func (p *PriceService) StartStream(ctx context.Context) {
	c := make(chan []byte)

	go func() {
		grpcServer := grpc.NewServer()
		pricesServer := server.NewPriceServer(ctx, c)
		protocol.RegisterPriceServiceServer(grpcServer, pricesServer)
	}()

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			prices, err := p.consumer.GetButchOfPrices(ctx)
			if err != nil {
				log.WithFields(log.Fields{
					"handler": "priceService",
					"action":  "get butch of prices",
				}).Errorf("unable to get butch of prices %v", err.Error())
			}

			for _, v := range prices {
				for _, v1 := range v.Messages {
					for _, v2 := range v1.Values {

						butchOfPrices := []models.Price{}

						if err = json.Unmarshal([]byte(v2.(string)), &butchOfPrices); err != nil {
							log.WithFields(log.Fields{
								"handler": "priceService",
								"action":  "unmarshal price",
							}).Errorf("unable to unmarshal price %v", err.Error())
						}

						go func() {
							c <- []byte(v2.(string))
						}()

						for _, price := range butchOfPrices {
							log.WithFields(log.Fields{
								"Price name ": price.Name,
								"Price bid ":  price.Bid,
								"Price ask":   price.Ask,
							}).Infof("Consumed price [%v]", price.Id)
						}
					}
				}
			}

			log.Infof("Tick at %v", t.UTC())
		}
	}
}

// NewPriceService creates new PriceService object
func NewPriceService(consumer *consumer.PriceConsumer) *PriceService {
	return &PriceService{consumer: consumer}
}
