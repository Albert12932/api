package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name" binding:"required"`
	Mail string `json:"mail" binding:"required"`
}

func main() {

	pool, err := pgxpool.New(context.Background(), "postgres://postgres:1234@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	r := gin.Default()
	err = r.SetTrustedProxies(nil)
	if err != nil {
		return
	}
	r.GET("/users/:id", getUserHandler(pool))
	r.POST("/users", createUserHandler(pool))
	r.DELETE("/users/delete/:id", deleteUserHandler(pool))
	r.PATCH("/users", patchUserHandler(pool))
	if err := r.Run(); err != nil {
		log.Println("Error while running")
		panic(err)
	}

}

func getUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err = pool.QueryRow(ctx, "SELECT id, name, mail from users where id = $1", id).Scan(&user.Id, &user.Name, &user.Mail)
		if err != nil {
			if err.Error() == "no rows in result set" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "Пользователь не найден",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка запроса getUser",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   user.Id,
			"name": user.Name,
			"mail": user.Mail,
		})
	}
}

func createUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser User

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Некорректные данные",
				"err":     err,
			})
			return
		}

		var UserID int

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := pool.QueryRow(ctx, "INSERT INTO users (name, mail) values ($1, $2) RETURNING id", newUser.Name, newUser.Mail).Scan(&UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Ошибка при добавлении пользователя",
				"message": err,
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Пользователь создан",
			"id":      UserID,
			"name":    newUser.Name,
			"mail":    newUser.Mail,
		})

	}
}

func deleteUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cmdTag, err := pool.Exec(ctx, "DELETE from users where id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка удаления пользователя",
			})
			return
		} else if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Такого пользователя нет",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Пользователь успешно удален",
			"id":      id,
		})

	}
}

func patchUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUserInfo User

		if err := c.ShouldBindJSON(&newUserInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Неверные данные",
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cmdTag, err := pool.Exec(ctx, "UPDATE users SET name=$1, mail=$2 where id=$3", newUserInfo.Name, newUserInfo.Mail, newUserInfo.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Ошибка обновления пользователя",
			})
			return
		} else if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Пользователя с таким id не существует",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Пользователь успешно обновлен",
			"id":      newUserInfo.Id,
			"name":    newUserInfo.Name,
			"mail":    newUserInfo.Mail,
		})
	}
}
