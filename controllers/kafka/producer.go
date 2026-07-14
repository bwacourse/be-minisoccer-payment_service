package kafka

import (
	configApp "payment-service/config"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Kafka struct {
	brokers []string
}

type IKafka interface {
	ProduceMessage(string, []byte) error
}

func NewKafkaProducer(brokers []string) IKafka {
	return &Kafka{
		brokers: brokers,
	}
}

func (k *Kafka) ProduceMessage(topic string, msg []byte) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = configApp.Config.Kafka.MaxRetry

	// init sarama producer
	producer, err := sarama.NewSyncProducer(k.brokers, config)

	if err != nil {
		logrus.Errorf("Error create producer: %v", err)
		return err
	}

	// make sure close producer when function return
	defer func(producer sarama.SyncProducer) {
		err = producer.Close()

		if err != nil {
			logrus.Errorf("Error close producer: %v", err)

			return
		}
	}(producer)

	//
	message := &sarama.ProducerMessage{
		Topic:   topic,
		Headers: nil,
		Value:   sarama.StringEncoder(msg),
	}

	partition, offset, err := producer.SendMessage(message)

	if err != nil {
		logrus.Errorf("Failed produce message to kafka topic %s: %v", topic, err)
		return err
	}

	logrus.Infof("Message produced successfully to topic %s at partition %d, offset %d", topic, partition, offset)

	return nil
}
