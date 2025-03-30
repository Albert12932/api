package main

import (
	"bdstudy/config"
	"bdstudy/routes"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

func main() {

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("error while connecting to db:%v", err)
	}
	defer db.Close()

	r := gin.Default()

	routes.RegisterUserRoutes(r, db)
	go func() {
		log.Println("🚀 Сервер запущен на порту 8080")
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("❌ Ошибка запуска сервера: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

}
