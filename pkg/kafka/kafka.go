package kafka

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/melwyn95/go-redis-kafka-demo/pkg/redis"
)

type KafkaProducer struct {
	producer *kafka.Producer
	Topic    *string
}

func NewProducer(host, port, topic string) (*KafkaProducer, error) {

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", host, port),
	})
	if err != nil {
		return nil, err
	}

	kp := KafkaProducer{
		producer: producer,
		Topic:    &topic,
	}

	return &kp, nil
}

func (kp *KafkaProducer) Produce(event []byte) error {
	err := kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: kp.Topic, Partition: kafka.PartitionAny},
		Value:          event,
	}, nil)

	if err != nil {
		return err
	}
	return nil
}
func (kp *KafkaProducer) Close() {
	fmt.Println("Closing kafka producer...")
	kp.producer.Close()
}

type KafkaConsumer struct {
	consumer *kafka.Consumer
}

func NewConsumer(host, port, topic, groupID string) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%s", host, port),
		"group.id":          groupID,
	})
	if err != nil {
		return nil, err
	}

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		return nil, err
	}

	kc := KafkaConsumer{
		consumer: consumer,
	}

	return &kc, nil
}

func (kc *KafkaConsumer) Close() {
	fmt.Println("Closing kafka consumer...")
	kc.consumer.Close()
}

type NumberEvent struct {
	EventType string `json:"type"`
	Number    int    `json:"number"`
}

func SpawnConsumer(kc *KafkaConsumer, redis *redis.Redis, quit <-chan os.Signal) {

	run := true
	for run == true {
		select {
		case <-quit:
			// Listen for os signals
			// close kakfa consumer
			kc.Close()
			run = false
		default:
			event := kc.consumer.Poll(100)

			if event == nil {
				continue
			}
			switch e := event.(type) {
			case kafka.AssignedPartitions:
				kc.consumer.Assign(e.Partitions)

			case kafka.RevokedPartitions:
				kc.consumer.Unassign()

			case *kafka.Message:
				var eve NumberEvent
				err := json.Unmarshal(e.Value, &eve)
				if err != nil {
					continue
				}
				if eve.Number%2 == 0 { // even
					redis.Incr("even")
				} else { //odd
					redis.Incr("odd")
				}
			}
		}
	}
}
