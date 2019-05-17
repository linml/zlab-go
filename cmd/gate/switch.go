package main

import (
	"fmt"
	"zlab/library/kafka/producer"
	"zlab/library/kafka/topic"
)

var (
	pd = producer.NewProducer()
)

//SwitchMessage 处理转发信息
func SwitchMessage(uid string, togame bool, data []byte) {

	if togame {
		pd.SwitchSend(topic.Gate2Game, uid, data)
	} else {
		pd.SwitchSend(topic.Gate2Cent, uid, data)
	}

	fmt.Println("handle msg, uid =", uid, "data length =", len(data))
}
