package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name" binding:"required"`
	Mail string `json:"mail" binding:"required"`
}

func main() {

	conn, err := pgx.Connect(context.Background(), "postgres://postgres:1234@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	var name string
	var mail string

	err = conn.QueryRow(context.Background(), "Select name, mail from users where id = 2").Scan(&name, &mail)
	if err != nil {
		log.Println(err, "Query failed")
		os.Exit(1)
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.GET("/users/:id", getUserHandler(conn))
	r.POST("/users/create", createUserHandler(conn))
	r.DELETE("/users/delete/:id", deleteUserHandler(conn))
	r.PATCH("/users/update", patchUserHandler(conn))
	r.Run()

}

func getUserHandler(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User

		id := c.Param("id")
		err := conn.QueryRow(context.Background(), "SELECT * from users where id = $1", id).Scan(&user.Id, &user.Name, &user.Mail)
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

func createUserHandler(conn *pgx.Conn) gin.HandlerFunc {
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

		err := conn.QueryRow(context.Background(), "INSERT INTO users (name, mail) values ($1, $2) RETURNING id", newUser.Name, newUser.Mail).Scan(&UserID)
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

func deleteUserHandler(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		cmdTag, err := conn.Exec(context.Background(), "DELETE from users where id = $1", id)
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

func patchUserHandler(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUserInfo User

		if err := c.ShouldBindJSON(&newUserInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Неверные данные",
			})
		}

		cmdTag, err := conn.Exec(context.Background(), "UPDATE users SET name=$1, mail=$2 where id=$3", newUserInfo.Name, newUserInfo.Mail, newUserInfo.Id)
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
