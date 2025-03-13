package kafk

//
//import (
//	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
//	"log"
//)
//
//func SendMessage(topic string, message string) error {
//	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
//	if err != nil {
//		return err
//	}
//	defer p.Close()
//
//	err = p.Produce(&kafka.Message{
//		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
//		Value:          []byte(message),
//	}, nil)
//
//	if err != nil {
//		return err
//	}
//
//	go func() {
//		for e := range p.Events() {
//			switch ev := e.(type) {
//			case *kafka.Message:
//				if ev.TopicPartition.Error != nil {
//					log.Printf("Ошибка отправки сообщения %v\n", ev.TopicPartition)
//				} else {
//					log.Printf("Сообщение отправлено в %v\n", ev.TopicPartition)
//				}
//			}
//		}
//	}()
//
//	p.Flush(15 * 1000)
//	return nil
//}
