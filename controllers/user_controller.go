package controllers

import (
	"bdstudy/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strconv"
	"time"
)

func GetUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err = pool.QueryRow(ctx, "SELECT id, name from users where id = $1", id).Scan(&user.Id, &user.Name)
		if err != nil {
			if err.Error() == "no rows in result set" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User is not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while getUser query",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id":   user.Id,
			"name": user.Name,
		})
	}
}

//func CreateUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		var newUser models.User
//
//		if err := c.ShouldBindJSON(&newUser); err != nil {
//			c.JSON(http.StatusBadRequest, gin.H{
//				"message": "Invalid data",
//				"err":     err,
//			})
//			return
//		}
//
//		var UserID int
//
//		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//		defer cancel()
//
//		err := pool.QueryRow(ctx, "INSERT INTO users (name) values ($1) RETURNING id", newUser.Name).Scan(&UserID)
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{
//				"error":   "Error while creating user",
//				"message": err.Error(),
//			})
//			return
//		}
//
//		c.JSON(http.StatusCreated, gin.H{
//			"message": "User successfully created",
//			"id":      UserID,
//			"name":    newUser.Name,
//		})
//
//	}
//}

func DeleteUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cmdTag, err := pool.Exec(ctx, "DELETE from users where id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while deleting user",
			})
			return
		} else if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User successfully deleted",
			"id":      id,
		})

	}
}

func PatchUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUserInfo models.User

		if err := c.ShouldBindJSON(&newUserInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid data",
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cmdTag, err := pool.Exec(ctx, "UPDATE users SET name=$1 where id=$2", newUserInfo.Name, newUserInfo.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while updating user",
			})
			return
		} else if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User successfully updated",
			"id":      newUserInfo.Id,
			"name":    newUserInfo.Name,
		})
	}
}

func RegisterUserHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser models.User

		if err := c.ShouldBindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Invalid data",
				"err":     err,
			})
			return
		}

		var UserID int

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := pool.QueryRow(ctx, "INSERT INTO users (name, password) values ($1, $2) RETURNING Id", newUser.Name, newUser.Password).Scan(&UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error while creating user",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "User successfully created",
			"id":      UserID,
			"name":    newUser.Name,
		})

	}
}

func GetLoginHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User

		username := c.Query("username")
		password := c.Query("password")

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := pool.QueryRow(ctx, "SELECT id, password FROM users WHERE name = $1", username).Scan(&user.Id, &user.Password)
		if err != nil {
			if err.Error() == "no rows in result set" {
				c.JSON(http.StatusNotFound, gin.H{
					"error": "User is not found",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while getUser query",
			})
			return
		}

		if password == user.Password {
			c.JSON(http.StatusOK, gin.H{
				"user_id": user.Id,
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Wrong password",
			})
			return
		}
	}
}
