package kafka

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type LogMessage struct {
	UserID    int    `json:"user_id"`
	Action    string `json:"action"`
	Details   string `json:"details"`
	Timestamp string `json:"timestamp"`
}

// KafkaLoggerProducer структура для логирования действий пользователей через Kafka
type KafkaLoggerProducer struct {
	Producer sarama.SyncProducer
}

// NewKafkaLoggerProducer создает новый Kafka Producer для логов
func NewKafkaLoggerProducer(brokers []string) (*KafkaLoggerProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Printf("❌ Ошибка создания Kafka Logger Producer: %v", err)
		return nil, err
	}

	log.Println("✅ Kafka Logger Producer успешно инициализирован")

	return &KafkaLoggerProducer{Producer: producer}, nil
}

// LogEvent отправляет событие логирования в Kafka
func (p *KafkaLoggerProducer) LogEvent(userID int, action string, details string) error {
	timestamp := time.Now().Format(time.RFC3339) // Время в ISO 8601 формате

	// Если details пустое, записываем "No details provided"
	if details == "" {
		details = "No details provided"
	}

	// Создаем объект JSON
	logMessage := LogMessage{
		Timestamp: timestamp,
		UserID:    userID,
		Action:    action,
		Details:   details,
	}

	// Кодируем в JSON
	jsonData, err := json.Marshal(logMessage)
	if err != nil {
		log.Printf("❌ Ошибка кодирования JSON: %v", err)
		return err
	}

	// Создаем сообщение для Kafka
	msg := &sarama.ProducerMessage{
		Topic: "user_logs",
		Value: sarama.ByteEncoder(jsonData),
	}

	partition, offset, err := p.Producer.SendMessage(msg)
	if err != nil {
		log.Printf("❌ Ошибка отправки лога в Kafka: %v", err)
		return err
	}

	log.Printf("✅ Лог отправлен в Kafka (partition: %d, offset: %d, time: %s, details: %s)\n",
		partition, offset, timestamp, details)
	return nil
}

// Close закрывает Producer
func (p *KafkaLoggerProducer) Close() {
	if err := p.Producer.Close(); err != nil {
		log.Printf("❌ Ошибка закрытия Kafka Logger Producer: %v", err)
	}
	log.Println("✅ Kafka Logger Producer закрыт")
}
