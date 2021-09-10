package server

import (
	"context"
	"github.com/INEFFABLE-games/PriceService/internal/protocol"
	log "github.com/sirupsen/logrus"
	"time"
)

type PriceServer struct {
	pricesChannels map[int]chan []byte
	ctx            context.Context

	protocol.UnimplementedPriceServiceServer
}

func (p *PriceServer) Send(stream protocol.PriceService_SendServer) error {

	currentIndex := len(p.pricesChannels) + 1
	currentChannel := make(chan []byte)

	log.WithFields(log.Fields{
		"handler ": "grpc send",
		"index ":   currentIndex,
	}).Info()

	p.pricesChannels[currentIndex] = currentChannel

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-p.ctx.Done():
			return nil
		case <-ticker.C:

			butchOfPrices := <-currentChannel
			if butchOfPrices == nil {
				continue
			}

			err := stream.Send(&protocol.SendReply{ButchOfPrices: butchOfPrices})
			if err != nil {
				log.WithFields(log.Fields{
					"handler ": "pricesServer(GRPC)",
					"action ":  "send request",
				}).Errorf("unable to send request %v", err.Error())

				delete(p.pricesChannels, currentIndex)
				return err
			}
		}
	}
}

func NewPriceServer(ctx context.Context, c map[int]chan []byte) *PriceServer {
	return &PriceServer{
		pricesChannels: c,
		ctx:            ctx,
	}
}
