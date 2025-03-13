package controllers

import (
	//"bdstudy/kafk"
	"bdstudy/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
	"strconv"
	"time"
)

func GetTasksHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while Atoi",
			})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

		defer cancel()

		rows, err := pool.Query(ctx, "SELECT id, userid, header, text, done from tasks where userid = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while getting tasks",
			})
			return
		}
		defer rows.Close()

		var tasks []models.Task

		for rows.Next() {
			var task models.Task
			if err := rows.Scan(&task.Id, &task.UserID, &task.Header, &task.Text, &task.Done); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Error while reading data from rows",
				})
				return
			}
			tasks = append(tasks, task)
		}
		if len(tasks) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Tasks not found",
			})
			return
		}
		c.JSON(http.StatusOK, tasks)
	}

}

func CreateTaskHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTask models.Task

		if err := c.ShouldBindJSON(&newTask); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid data format",
				"message": err.Error(),
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel()

		_, err := pool.Exec(ctx, "Insert into tasks (userId, header, text, done) values ($1, $2, $3, $4)", newTask.UserID, newTask.Header, newTask.Text, newTask.Done)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Error while creating task",
				"message": err,
			})
			return
		}

		//err = kafk.SendMessage("tasks", "Task created with ID:"+strconv.Itoa(newTask.Id))
		//if err != nil {
		//	log.Println("Error while sending message with kafk")
		//}

		c.JSON(http.StatusOK, gin.H{
			"message": "Task successfully created",
			"taskId":  newTask.Id,
		})

	}
}

func DeleteTaskHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cmdTag, err := pool.Exec(ctx, "DELETE from tasks where id = $1", id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while deleting task",
			})
			return
		} else if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Task successfully deleted",
			"id":      id,
		})

	}
}

func PatchTaskHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newTaskInfo models.Task

		if err := c.ShouldBindJSON(&newTaskInfo); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid data",
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		cmdTag, err := pool.Exec(ctx, "UPDATE tasks SET userId=$1, header=$2, text=$3, done=$4 where id=$5", newTaskInfo.UserID, newTaskInfo.Header, newTaskInfo.Text, newTaskInfo.Done, newTaskInfo.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error while updating task",
			})
			return
		} else if cmdTag.RowsAffected() == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Task not found",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Task successfully updated",
			"id":      newTaskInfo.Id,
			"userId":  newTaskInfo.UserID,
			"header":  newTaskInfo.Header,
			"text":    newTaskInfo.Text,
			"done":    newTaskInfo.Done,
		})
	}
}

func SwitchTaskHandler(pool *pgxpool.Pool) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid id",
			})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		defer cancel()

		var newState bool

		err = pool.QueryRow(ctx, "update tasks set done= NOT done where id = $1 RETURNING done", id).Scan(&newState)
		if err != nil {
			if err.Error() == "no rows in result set" {
				c.JSON(http.StatusNotFound, gin.H{"error": "Задача не найдена"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "error while updating task state",
				"message": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Task state updated successfully",
			"state":   newState,
		})

	}
}
