package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type Task struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
	DueDate  string `json:"dueDate"`
}

var tasks = []Task{
	{ID: 1, Name: "Task 1", Category: "Category 1", DueDate: "2023-01-01"},
	{ID: 2, Name: "Task 2", Category: "Category 2", DueDate: "2023-03-15"},
	{ID: 3, Name: "Task 3", Category: "Category 3", DueDate: "2023-06-30"},
}

func main() {
	e := echo.New()

	e.GET("/", getIntro)
	e.GET("/tasks", getTasks)
	e.GET("/tasks/:id", getTaskByID)
	e.POST("/tasks", createTask)
	e.PUT("/tasks/:id", updateTask)
	e.DELETE("/tasks/:id", deleteTask)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := e.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}

func getTasks(c echo.Context) error {
	category := c.QueryParam("category")
	return c.JSON(http.StatusOK, filterTasks(category))
}

func filterTasks(category string) []Task {
	result := make([]Task, 0)
	for _, task := range tasks {
		if category != "" && task.Category != category {
			continue
		}
		result = append(result, task)
	}
	return result
}

func getTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid task ID"})
	}

	for _, task := range tasks {
		if task.ID == id {
			return c.JSON(http.StatusOK, task)
		}
	}

	return c.JSON(http.StatusNotFound, echo.Map{"error": "Task not found"})
}

func createTask(c echo.Context) error {
	var task Task
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	if !validateDueDate(task.DueDate) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid due date (format: YYYY-MM-DD)"})
	}

	task.ID = len(tasks) + 1
	tasks = append(tasks, task)

	return c.JSON(http.StatusCreated, task)
}

func updateTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid task ID"})
	}

	var updatedTask Task
	if err := c.Bind(&updatedTask); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request payload"})
	}

	if !validateDueDate(updatedTask.DueDate) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid due date (format: YYYY-MM-DD)"})
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Name = updatedTask.Name
			tasks[i].Category = updatedTask.Category
			tasks[i].DueDate = updatedTask.DueDate
			return c.JSON(http.StatusOK, tasks[i])
		}
	}

	return c.JSON(http.StatusNotFound, echo.Map{"error": "Task not found"})
}

func deleteTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid task ID"})
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("Task %d deleted", id)})
		}
	}

	return c.JSON(http.StatusNotFound, echo.Map{"error": "Task not found"})
}

func validateDueDate(date string) bool {
	regex := `^\d{4}-\d{2}-\d{2}$`
	match, _ := regexp.MatchString(regex, date)
	return match
}

func getIntro(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the task manager!")
}
