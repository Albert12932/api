package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/IBM/sarama"
)

// LogEntry структура для логов
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	UserID    int    `json:"user_id"`
	Action    string `json:"action"`
	Details   string `json:"details"`
}

// KafkaLogConsumer читает логи из Kafka
type KafkaLogConsumer struct {
	Consumer sarama.ConsumerGroup
	Topic    string
}

// NewKafkaLogConsumer создает новый Kafka Consumer
func NewKafkaLogConsumer(brokers []string, group, topic string) (*KafkaLogConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V3_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}

	consumer, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Printf("❌ Ошибка создания Kafka Log Consumer: %v", err)
		return nil, err
	}

	log.Println("✅ Kafka Log Consumer успешно инициализирован")

	return &KafkaLogConsumer{
		Consumer: consumer,
		Topic:    topic,
	}, nil
}

// ConsumeLogs читает логи из Kafka и записывает их в `logs.json`
func (kc *KafkaLogConsumer) ConsumeLogs(ctx context.Context) {
	file, err := os.OpenFile("logs.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("❌ Ошибка открытия файла логов: %v", err)
	}
	defer file.Close()

	handler := &logConsumerHandler{file}

	for {
		select {
		case <-ctx.Done(): // Если контекст завершён (Ctrl+C), выходим из цикла
			log.Println("🛑 Остановка Kafka Log Consumer...")
			return
		default:
			err := kc.Consumer.Consume(ctx, []string{kc.Topic}, handler)
			if err != nil {
				// Если ошибка связана с закрытием контекста, выходим без ошибки
				if err == context.Canceled || err.Error() == "kafka: tried to use a consumer group that was closed" {
					log.Println("✅ Kafka Consumer корректно завершён")
					return
				}
				log.Printf("❌ Ошибка чтения логов из Kafka: %v\n", err)
			}
		}
	}
}

// Close закрывает Consumer
func (kc *KafkaLogConsumer) Close() {
	if err := kc.Consumer.Close(); err != nil {
		log.Printf("❌ Ошибка закрытия Kafka Log Consumer: %v\n", err)
	}
	log.Println("✅ Kafka Log Consumer закрыт")
}

// logConsumerHandler обрабатывает сообщения из Kafka
type logConsumerHandler struct {
	File *os.File
}

func (h logConsumerHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h logConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (h logConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var logEntry LogEntry
		err := json.Unmarshal(message.Value, &logEntry)
		if err != nil {
			log.Printf("❌ Ошибка парсинга JSON: %v\n JSON: %s\n", err, string(message.Value))
			continue
		}

		// Если details пустое, заменяем его
		if logEntry.Details == "" {
			logEntry.Details = "No details provided"
		}

		// Записываем JSON в файл
		jsonData, err := json.Marshal(logEntry)
		if err != nil {
			log.Printf("❌ Ошибка кодирования JSON: %v\n", err)
			continue
		}

		_, err = h.File.Write(jsonData)
		if err != nil {
			log.Printf("❌ Ошибка записи в файл: %v\n", err)
		}

		// Добавляем перенос строки для читабельности
		h.File.WriteString("\n")

		// Логируем в консоль
		log.Printf("📩 Лог записан: %s\n", string(jsonData))

		session.MarkMessage(message, "")
	}
	return nil
}
