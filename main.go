package main

import (
	"bdstudy/config"
	"bdstudy/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("error while connecting to db")
	}
	defer db.Close()

	r := gin.Default()

	routes.RegisterTaskRoutes(r, db)
	routes.RegisterUserRoutes(r, db)

	log.Println("Server started at port 8080")
	err = r.Run(":8080")
	if err != nil {
		log.Fatal("error while running")
	}

}
