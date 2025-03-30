package routes

import (
	"bdstudy/controllers"
	"bdstudy/kafka"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterTaskRoutes(r *gin.Engine, db *pgxpool.Pool) {
	r.GET("/tasks/:id", controllers.GetTasksHandler(db))
	r.POST("/tasks", controllers.CreateTaskHandler(db))
	r.DELETE("/tasks/delete/:id", controllers.DeleteTaskHandler(db))
	r.PATCH("/tasks", controllers.PatchTaskHandler(db))
	r.PATCH("/tasks/:id", controllers.SwitchTaskHandler(db))
}
