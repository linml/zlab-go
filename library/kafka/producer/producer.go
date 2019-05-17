package producer

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"zlab/library/kafka/topic"

	"github.com/Shopify/sarama"
)

var (
	addr      = flag.String("addr", ":8080", "The address to bind to")
	brokers   = flag.String("brokers", os.Getenv("KAFKA_PEERS"), "The Kafka brokers to connect to, as a comma separated list")
	verbose   = flag.Bool("verbose", false, "Turn on Sarama logging")
	certFile  = flag.String("certificate", "", "The optional certificate file for client authentication")
	keyFile   = flag.String("key", "", "The optional key file for client authentication")
	caFile    = flag.String("ca", "", "The optional certificate authority file for TLS client authentication")
	verifySsl = flag.Bool("verify", false, "Optional verify ssl certificates chain")
)

//NewProducer ..
func NewProducer() *Producer {
	flag.Parse()

	if *verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	if *brokers == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	brokerList := strings.Split(*brokers, ",")
	log.Printf("Kafka brokers: %s", strings.Join(brokerList, ", "))
	return &Producer{
		Sync:  newDataCollector(brokerList),
		Async: newAccessLogProducer(brokerList),
	}

}

//Producer 队列生产者
type Producer struct {
	Sync  sarama.SyncProducer
	Async sarama.AsyncProducer
}

//Send 转发消息 protobuf
func (s *Producer) SwitchSend(topic, uid string, data []byte) error {
	var msg = &sarama.ProducerMessage{
		Value: sarama.ByteEncoder(data),  //数据
		Key:   sarama.StringEncoder(uid), //用户id
		Topic: topic,                     //game gate cent
	}
	p, off, err := s.Sync.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg error", err)
		return err
	}
	fmt.Println("send msg on partion", p, "offset", off)
	return nil
}

//SwitchSendAsync 异步发送消息 protobuf
func (s *Producer) SwitchSendAsync(topic, uid string, data []byte) {
	var msg = &sarama.ProducerMessage{
		Value: sarama.ByteEncoder(data),  //数据
		Key:   sarama.StringEncoder(uid), //用户id
		Topic: topic,                     //game gate cent
	}

	s.Async.Input() <- msg
}

//LogAsync 异步log
func (s *Producer) LogAsync(text string) {
	var msg = &sarama.ProducerMessage{
		Value: sarama.StringEncoder(text),
		Topic: topic.Log,
	}
	s.Async.Input() <- msg
}

//DbAsync 数据库写
func (s *Producer) DbAsync(cmd []byte) {
	var msg = &sarama.ProducerMessage{
		Value: sarama.ByteEncoder(cmd),
		Topic: topic.DB,
	}
	s.Async.Input() <- msg
}

func newDataCollector(brokerList []string) sarama.SyncProducer {

	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true
	// tlsConfig := createTlsConfiguration()
	// if tlsConfig != nil {
	// 	config.Net.TLS.Config = tlsConfig
	// 	config.Net.TLS.Enable = true
	// }

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

func newAccessLogProducer(brokerList []string) sarama.AsyncProducer {

	// For the access log, we are looking for AP semantics, with high throughput.
	// By creating batches of compressed messages, we reduce network I/O at a cost of more latency.
	config := sarama.NewConfig()
	// tlsConfig := createTlsConfiguration()
	// if tlsConfig != nil {
	// 	config.Net.TLS.Enable = true
	// 	config.Net.TLS.Config = tlsConfig
	// }
	config.Producer.RequiredAcks = sarama.WaitForLocal       // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy   // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond // Flush batches every 500ms

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write access log entry:", err)
		}
	}()

	return producer
}
