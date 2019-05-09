package main

import (
	"fmt"
	"reflect"
	"sync"
	"zlab/library/reque"

	"github.com/golang/protobuf/proto"

	"zlab/pb"

	"github.com/go-redis/redis"
)

func main() {
	var wg = sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		reque.Register(&pb.HelloRequest{}, sayHello)
		reque.Init()
	}()

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	go func() {
		var msg = &pb.HelloRequest{
			Greeting: "hello m1",
		}

		var t = reflect.TypeOf(msg)
		var channel = t.String()

		buf, err := proto.Marshal(msg)
		if err != nil {
			fmt.Println("marshal error", err)
		}

		fmt.Println(channel, msg)
		err = client.Publish(channel, buf).Err()
		if err != nil {
			fmt.Println("m1 error", err)
		}

		err = client.Publish(channel, buf).Err()
		if err != nil {
			fmt.Println("m2 error", err)
		}
	}()

	wg.Wait()
}

func sayHello(msg proto.Message) {
	SayHello(msg.(*pb.HelloRequest))
}

func SayHello(msg *pb.HelloRequest) {
	fmt.Println("HelloRequest invoke", msg)
}
