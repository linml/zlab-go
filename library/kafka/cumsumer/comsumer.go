package comsumer

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"zlab/library/kafka/cons"
	"zlab/util"

	"github.com/golang/protobuf/proto"

	"github.com/Shopify/sarama"
)

// Sarma configuration options
var (
	brokers = "localhost:9020"
	verbose = false
)

//Init 初始化
func Init(group string, topics ...string) {
	log.Println("Starting a new Sarama consumer")

	if verbose {
		sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()

	/**
	 * Setup a new Sarama consumer group
	 */
	consumer := Consumer{}

	ctx := context.Background()
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			consumer.ready = make(chan bool, 0)
			err := client.Consume(ctx, topics, &consumer)
			if err != nil {
				panic(err)
			}
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	<-sigterm // Await a sigterm signal before safely closing the consumer

	err = client.Close()
	if err != nil {
		panic(err)
	}
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready chan bool
	m     map[string]Message
}

type Message struct {
	msg    proto.Message
	single bool
	fun    func(msg proto.Message)
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(consumer.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {

		switch message.Topic {
		case "db":
			fmt.Println("mongodb udpate")
		case "gate", "game", "cent":
			var stype, id string
			if message.Headers != nil && len(message.Headers) > 0 {
				for _, header := range message.Headers {
					switch util.ToString(header.Key) {
					case cons.TYPEKEY:
						stype = util.ToString(header.Value)
					case cons.IDKEY:
						id = util.ToString(header.Value)
					}

				}
			}
			if stype == "sys" {

			}
			fmt.Println(stype, id)
			if m, ok := consumer.m[util.ToString(message.Key)]; ok {
				err := proto.Unmarshal(message.Value, m.msg)
				if err != nil {
					fmt.Println("unmarshal msg error")
				} else {
					m.fun(m.msg)
				}

				if m.single {
					session.MarkMessage(message, "")
				}
			} else {
				fmt.Println("unkown msg", message.Value)
				session.MarkMessage(message, "")
			}
		case "log":
			var text = util.ToString(message.Value)
			fmt.Println("log", text)
			session.MarkMessage(message, "")
		}
	}

	return nil
}
