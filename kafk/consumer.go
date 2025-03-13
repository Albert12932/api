package kafk

//
//import (
//	"fmt"
//	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
//	"log"
//)
//
//func StartConsumer() {
//	c, err := kafka.NewConsumer(&kafka.ConfigMap{
//		"bootstrap.servers": "localhost:9092",
//		"group.id":          "task-consumer-group",
//		"auto.offset.reset": "earliest",
//	})
//	if err != nil {
//		log.Fatalf("Ошибка создания kafka Consumer %v", err)
//	}
//	defer c.Close()
//
//	topic := "tasks"
//	c.SubscribeTopics([]string{topic}, nil)
//
//	for {
//		msg, err := c.ReadMessage(-1)
//		if err == nil {
//			fmt.Printf("Получено сообщение %s\n", string(msg.Value))
//		} else {
//			fmt.Printf("Ошибка при чтении сообщения %v\n", err)
//		}
//	}
//}
