package message

import (
	"zlab/library/kafka/producer"
)

//Producer 生产者
var Producer *producer.Producer

func init() {
	Producer = producer.NewProducer()
}
