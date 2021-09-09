package server

import (
	"context"
	"github.com/INEFFABLE-games/PriceService/internal/protocol"
	log "github.com/sirupsen/logrus"
	"time"
)

type PriceServer struct {
	pricesChannel chan []byte
	ctx           context.Context

	protocol.UnimplementedPriceServiceServer
}

func (p *PriceServer) Send(stream protocol.PriceService_SendServer) error {

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-p.ctx.Done():
			return nil
		case <-ticker.C:

			butchOfPrices := <-p.pricesChannel
			if butchOfPrices == nil {
				continue
			}

			log.Infof("get new prices %v", butchOfPrices)

			err := stream.Send(&protocol.SendReply{ButchOfPrices: butchOfPrices})
			if err != nil {
				log.WithFields(log.Fields{
					"handler ": "pricesServer(GRPC)",
					"action ":  "send request",
				}).Errorf("unable to send request %v", err.Error())
			}
		}
	}
}

func NewPriceServer(ctx context.Context, c chan []byte) *PriceServer {
	return &PriceServer{
		pricesChannel: c,
		ctx:           ctx,
	}
}
