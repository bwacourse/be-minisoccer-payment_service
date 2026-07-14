package kafka

type KafkaRegistry struct {
	brokers []string
}

type IKafkaRegistry interface {
	GetKafkaProducer() IKafka
}

func NewKafkaRegistry(brokers []string) IKafkaRegistry {
	return &KafkaRegistry{
		brokers: brokers,
	}
}

// GetKafkaProducer implements [IKafkaRegistry].
func (k *KafkaRegistry) GetKafkaProducer() IKafka {
	return NewKafkaProducer(k.brokers)
}
