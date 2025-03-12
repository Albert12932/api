package routes

import (
	"bdstudy/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRoutes(r *gin.Engine, db *pgxpool.Pool) {
	r.GET("/users/:id", controllers.GetUserHandler(db))
	r.POST("/users", controllers.CreateUserHandler(db))
	r.DELETE("/users/delete/:id", controllers.DeleteUserHandler(db))
	r.PATCH("/users", controllers.PatchUserHandler(db))
}
