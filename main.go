package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"time"
)

func getRedis() *redis.Client {
	var (
		RedisAddres    = "redis-14436.c54.ap-northeast-1-2.ec2.cloud.redislabs.com:14436"
		RedisPassword  = "Drektarov3698!"
		ReddisUserName = "Admin"
	)

	client := redis.NewClient(&redis.Options{
		Addr:     RedisAddres,
		Password: RedisPassword,
		Username: ReddisUserName,
		DB:       0,
	})

	res, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)

	return client
}

func main(){

	client := getRedis()

	ctx := context.Background()

	for {
		price, err := client.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"Prices","0"},
			Count:   0,
			Block:   2 * time.Second,
		}).Result()

		if price == nil{
			continue
		}

		if err != nil {
			log.Println(err)
		} else {
			for _, val := range price {
				for k, v := range val.Messages {
					log.Printf("[%d]: %v \n", k, v)
				}

				fmt.Println("")
			}
		}

		err = client.Do(ctx, "DEL", "Prices").Err();if err != nil{
			log.WithFields(log.Fields{
				"handler" : "main",
				"action" : "clear redis db",
			}).Errorf("unable to clear redis db %v",err.Error())
		}

		time.Sleep(1 * time.Second)
	}

}
