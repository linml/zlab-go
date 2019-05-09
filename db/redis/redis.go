package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

var client *redis.Client

func init() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pb := client.Subscribe("s1")
	var ch = pb.Channel()
	for {
		select {
		case v := <-ch:
			fmt.Println(v.String())
		}
	}
	fmt.Println("redis is init")
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func Redis() *redis.Client {
	return client
}
