package producer

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
	"zlab/library/kafka/cons"
	"zlab/util"

	"github.com/golang/protobuf/proto"

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
		DataCollector:     newDataCollector(brokerList),
		AccessLogProducer: newAccessLogProducer(brokerList),
	}

}

//Producer 队列生产者
type Producer struct {
	DataCollector     sarama.SyncProducer
	AccessLogProducer sarama.AsyncProducer
}

//Send 发送消息 protobuf
func (s *Producer) Send(topic, stype, id string, pb proto.Message) error {
	var buf, err = proto.Marshal(pb)
	if err != nil {
		fmt.Println("marshal msg error", err)
		return err
	}
	var msg = &sarama.ProducerMessage{
		Value: sarama.ByteEncoder(buf),
		Key:   sarama.StringEncoder(reflect.TypeOf(pb).String()),
		Topic: topic,
		Headers: []sarama.RecordHeader{
			sarama.RecordHeader{
				Key:   util.ToBytes("type"),
				Value: util.ToBytes(stype),
			},
			sarama.RecordHeader{
				Key:   util.ToBytes("id"),
				Value: util.ToBytes(id),
			},
		},
	}
	p, off, err := s.DataCollector.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg error", err)
		return err
	}
	fmt.Println("send msg on partion", p, "offset", off)
	return nil
}

//SendAsync 异步发送消息 protobuf
func (s *Producer) SendAsync(topic, stype, id string, pb proto.Message) {
	var buf, err = proto.Marshal(pb)
	if err != nil {
		fmt.Println("marshal msg error", err)
		return
	}

	var msg = &sarama.ProducerMessage{
		Value: sarama.ByteEncoder(buf),
		Key:   sarama.StringEncoder(reflect.TypeOf(pb).String()),
		Topic: topic,
		Headers: []sarama.RecordHeader{
			sarama.RecordHeader{
				Key:   util.ToBytes(cons.TYPEKEY),
				Value: util.ToBytes(stype),
			},
			sarama.RecordHeader{
				Key:   util.ToBytes(cons.IDKEY),
				Value: util.ToBytes(id),
			},
		},
	}

	s.AccessLogProducer.Input() <- msg
}

//LogAsync 异步log
func (s *Producer) LogAsync(topic, text string) {

	var msg = &sarama.ProducerMessage{
		Value: sarama.StringEncoder(text),
		Topic: topic,
	}
	s.AccessLogProducer.Input() <- msg
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
