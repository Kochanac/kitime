package kafka

import (
	"context"
	"fmt"
	"github.com/Kochanac/kitime/service/pkg/config"
	"github.com/Shopify/sarama"
	"log"
	"os"
)

type Producer interface {
	Produce(ctx context.Context, message string) error
}

type SaramaProducer struct {
	sarama.SyncProducer
}

func InitProducer(c config.Config) (Producer, error) {
	sarama.Logger = log.New(os.Stdout, "", log.Ltime)

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Retry.Max = 5
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.ClientID = "vobla-producer"

	prd, err := sarama.NewSyncProducer([]string{c.GetKafkaHost()}, saramaConfig)
	if err != nil {
		return nil, err
	}
	return &SaramaProducer{prd}, err
}

func (producer *SaramaProducer) Produce(ctx context.Context, message string) error {
	c := ctx.Value("config").(config.Config)

	msg := &sarama.ProducerMessage{
		Topic: c.GetKafkaTopic(),
		Value: sarama.StringEncoder(message),
	}
	p, o, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Println("Error publish: ", err.Error())
		return err
	}

	fmt.Println("Partition: ", p)
	fmt.Println("Offset: ", o)
	return nil
}
