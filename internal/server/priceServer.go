package server

import (
	"context"
	protocol2 "github.com/INEFFABLE-games/PriceService/protocol"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"time"
)

type PriceServer struct {
	pricesChannels map[string]chan []byte
	ctx            context.Context

	protocol2.UnimplementedPriceServiceServer
}

func (p *PriceServer) Send(stream protocol2.PriceService_SendServer) error {

	currentIndex := uuid.NewV4().String()

	currentChannel := make(chan []byte)

	log.WithFields(log.Fields{
		"handler ": "grpc send",
		"index ":   currentIndex,
	}).Info()

	p.pricesChannels[currentIndex] = currentChannel

	ticker := time.NewTicker(1 * time.Millisecond)
	for {
		select {
		case <-p.ctx.Done():
			delete(p.pricesChannels, currentIndex)
			return nil
		case <-ticker.C:

			if p.pricesChannels[currentIndex] == nil {
				continue
			}

			butchOfPrices := <-p.pricesChannels[currentIndex]

			err := stream.Send(&protocol2.SendReply{ButchOfPrices: butchOfPrices})
			if err != nil {
				log.WithFields(log.Fields{
					"handler ": "pricesServer(GRPC)",
					"action ":  "send reply",
				}).Errorf("unable to send reply %v", err.Error())

				delete(p.pricesChannels, currentIndex)
				return err
			}
		}
	}
}

func NewPriceServer(ctx context.Context, c map[string]chan []byte) *PriceServer {
	return &PriceServer{
		pricesChannels: c,
		ctx:            ctx,
	}
}
