package routes

import (
	"bdstudy/controllers"
	"bdstudy/kafka"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterTaskRoutes(r *gin.Engine, db *pgxpool.Pool, logger *kafka.KafkaLoggerProducer) {
	r.GET("/tasks/:id", controllers.GetTasksHandler(db, logger))
	r.POST("/tasks", controllers.CreateTaskHandler(db, logger))
	r.DELETE("/tasks/delete/:id", controllers.DeleteTaskHandler(db, logger))
	r.PATCH("/tasks", controllers.PatchTaskHandler(db, logger))
	r.PATCH("/tasks/:id", controllers.SwitchTaskHandler(db, logger))
}
