package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID       int        `json:"id"`
	Name     string     `json:"name"`
	Category string     `json:"category"`
	DueDate  string     `json:"dueDate"`
	Status   TaskStatus `json:"status"`
}

type TaskUpdate struct {
	Name     *string     `json:"name"`
	Category *string     `json:"category"`
	DueDate  *string     `json:"dueDate"`
	Status   *TaskStatus `json:"status"`
}

var tasks = []Task{
	{ID: 1, Name: "Task 1", Category: "Category 1", DueDate: "2023-01-01", Status: Todo},
	{ID: 2, Name: "Task 2", Category: "Category 2", DueDate: "2023-03-15", Status: Doing},
	{ID: 3, Name: "Task 3A", Category: "Category 3", DueDate: "2023-06-30", Status: Doing},
	{ID: 3, Name: "Task 3B", Category: "Category 3", DueDate: "2024-01-31", Status: Done},
}

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

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
	filteredTasks := filterTasks(category)
	sortBy := c.QueryParam("sortBy")
	sortDir := c.QueryParam("sortDir")
	sortedFilteredTasks := sortTasks(filteredTasks, sortBy, sortDir)
	return c.JSON(http.StatusOK, sortedFilteredTasks)
}

func sortTasks(taskList []Task, by string, dir string) []Task {
	fmt.Printf("Sorting tasks by %s, direction %s\n", by, dir)
	dir = strings.ToLower(dir)
	t := make([]Task, len(taskList))
	copy(t, taskList)
	sort.Slice(t, func(i, j int) bool {
		switch by {
		case "id":
			if dir == "desc" {
				return taskList[i].ID > taskList[j].ID
			}
			return taskList[i].ID < taskList[j].ID
		case "name":
			if dir == "desc" {
				return taskList[i].Name > taskList[j].Name
			}
			return taskList[i].Name < taskList[j].Name
		case "category":
			if dir == "desc" {
				return taskList[i].Category > taskList[j].Category
			}
			return taskList[i].Category < taskList[j].Category
		case "dueDate":
			if dir == "desc" {
				return taskList[i].DueDate > taskList[j].DueDate
			}
			return taskList[i].DueDate < taskList[j].DueDate
		case "status":
			if dir == "desc" {
				return taskList[i].DueDate > taskList[j].DueDate
			}
			return taskList[i].DueDate < taskList[j].DueDate
		default:
			// Sort by id asc by default
			if dir == "desc" {
				return taskList[i].ID > taskList[j].ID
			}
			return taskList[i].ID < taskList[j].ID
		}
	})
	return t
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("Invalid task ID: %s", err)})
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("Invalid request payload: %s", err)})
	}

	if !isValidDate(task.DueDate) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid due date (format: YYYY-MM-DD)"})
	}

	task.ID = len(tasks) + 1
	tasks = append(tasks, task)

	return c.JSON(http.StatusCreated, task)
}

func updateTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("Invalid task ID: %s", err)})
	}

	var taskUpdate TaskUpdate
	if err := c.Bind(&taskUpdate); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("Invalid request payload: %s", err)})
	}

	if taskUpdate.DueDate != nil && !isValidDate(*taskUpdate.DueDate) {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid due date (format: YYYY-MM-DD)"})
	}

	for i := range tasks {
		if tasks[i].ID == id {
			if taskUpdate.Name != nil {
				tasks[i].Name = *taskUpdate.Name
			}
			if taskUpdate.Category != nil {
				tasks[i].Category = *taskUpdate.Category
			}
			if taskUpdate.DueDate != nil {
				tasks[i].DueDate = *taskUpdate.DueDate
			}
			return c.JSON(http.StatusOK, tasks[i])
		}
	}

	return c.JSON(http.StatusNotFound, echo.Map{"error": "Task not found"})
}

func deleteTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": fmt.Sprintf("Invalid task ID: %s", err)})
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return c.JSON(http.StatusOK, echo.Map{"message": fmt.Sprintf("Task %d deleted", id)})
		}
	}

	return c.JSON(http.StatusNotFound, echo.Map{"error": "Task not found"})
}

func isValidDate(dateString string) bool {
	layout := "2006-01-02"
	_, err := time.Parse(layout, dateString)
	return err == nil
}

func getIntro(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the task manager!")
}
