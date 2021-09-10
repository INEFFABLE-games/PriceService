package service

import (
	"context"
	"encoding/json"
	"github.com/INEFFABLE-games/PriceService/internal/consumer"
	"github.com/INEFFABLE-games/PriceService/models"
	"github.com/INEFFABLE-games/PriceService/protocol"
	log "github.com/sirupsen/logrus"
	"time"
)

// PriceService is structure for PriceService object
type PriceService struct {
	consumer *consumer.PriceConsumer
	server   *protocol.PriceServiceServer
}

// StartStream starts infinity cycle to get butch of prices from redis
func (p *PriceService) StartStream(ctx context.Context, channels map[int]chan []byte) {

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
							for _, v := range channels {
								v <- []byte(v2.(string))
							}
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
