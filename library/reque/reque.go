package reque

import (
	"fmt"
	"reflect"
	"zlab/util"

	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
)

type Message struct {
	msg proto.Message
	fun func(proto.Message)
}

var m = make(map[string]Message, 100)

func Register(pb proto.Message, fun func(proto.Message)) {
	var t = reflect.TypeOf(pb)
	m[t.String()] = Message{
		msg: pb,
		fun: fun,
	}
	fmt.Println("register", t.String())
}

func Init() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	var channels = make([]string, 0, len(m))
	for k := range m {
		channels = append(channels, k)
	}
	var ps = client.Subscribe(channels...)
	fmt.Println("sub channels", channels)
	ch := ps.Channel()
	for {
		select {
		case msg := <-ch:
			if m, ok := m[msg.Channel]; ok {
				err := proto.Unmarshal(util.ToBytes(msg.Payload), m.msg)
				if err != nil {
					fmt.Println("unmarshal error", err)
				} else {
					m.fun(m.msg)
				}
			} else {
				fmt.Println("unkown msg", msg.Channel, msg.Payload)
			}
		}
	}
}
