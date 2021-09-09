package consumer

import (
	"context"
	"github.com/INEFFABLE-games/PriceService/internal/config"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
)

// PriceConsumer structure for PriceConsumer object
type PriceConsumer struct {
	consumer *redis.Client
}

// GetButchOfPrices get and return butch of prices from redis
func (p PriceConsumer) GetButchOfPrices(ctx context.Context) ([]redis.XStream, error) {
	price, err := p.consumer.XRead(ctx, &redis.XReadArgs{
		Streams: []string{"Prices", "0"},
		Count:   0,
		Block:   2 * time.Second,
	}).Result()

	if err != nil {
		log.WithFields(log.Fields{
			"handler ": "PriceConsumer",
			"action ":  "read stream",
		}).Errorf("unable to read stream %v", err.Error())

		return nil, err
	}
	err = p.consumer.Do(ctx, "DEL", "Prices").Err()
	if err != nil {
		log.WithFields(log.Fields{
			"handler": "PriceConsumer",
			"action":  "clear redis db",
		}).Errorf("unable to clear redis db %v", err.Error())
	}
	return price, err

}

// getRedis returns new redis client connection
func getRedis(cfg *config.Config) *redis.Client {
	var (
		RedisAddres   = cfg.RedisAddress
		RedisPassword = cfg.RedisPassword
		RedisUserName = cfg.RedisUserName
	)

	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddres,
		Password: RedisPassword,
		Username: RedisUserName,
		DB:       0,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Errorf("unable to ping redis connection %v", err.Error())
	}

	return client
}

// NewPriceConsumer creates new PriceConsumer object
func NewPriceConsumer(cfg *config.Config) *PriceConsumer {
	consumer := getRedis(cfg)
	return &PriceConsumer{consumer: consumer}
}
