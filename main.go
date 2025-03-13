package main

import (
	"bdstudy/config"
	"bdstudy/kafka"
	"bdstudy/routes"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	logger := log.New(os.Stdout, "[TODO-KAFKA] ", log.LstdFlags)

	// Адрес Kafka-брокера
	brokers := []string{"localhost:9092"}

	// Инициализация Kafka Logger Producer
	logProducer, err := kafka.NewKafkaLoggerProducer(brokers)
	if err != nil {
		logger.Fatalf("❌ Ошибка создания Kafka Logger Producer: %v", err)
	}
	defer logProducer.Close()

	// Запуск Kafka Log Consumer (для записи логов в файл)
	ctx, cancel := context.WithCancel(context.Background())
	logConsumer, err := kafka.NewKafkaLogConsumer(brokers, "log-group", "user_logs")
	if err != nil {
		logger.Fatalf("❌ Ошибка создания Kafka Log Consumer: %v", err)
	}
	go logConsumer.ConsumeLogs(ctx)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("error while connecting to db")
	}
	defer db.Close()

	r := gin.Default()

	routes.RegisterTaskRoutes(r, db, logProducer)
	routes.RegisterUserRoutes(r, db)

	go func() {
		log.Println("🚀 Сервер запущен на порту 8080")
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("❌ Ошибка запуска сервера: %v", err)
		}
	}()

	<-stop
	logger.Println("🛑 Завершение работы сервера...")

	logConsumer.Close()

	cancel()

	time.Sleep(1 * time.Second)

}
