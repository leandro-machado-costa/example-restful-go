package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "Learn Go", Completed: false},
	{ID: "2", Item: "Build RESTful API", Completed: false},
	{ID: "3", Item: "Learn Docker", Completed: false},
}

func getTodos(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, todos)
}

func addTodo(c *gin.Context) {
	var newTodo todo
	if err := c.BindJSON(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todos = append(todos, newTodo)
	c.IndentedJSON(http.StatusCreated, newTodo)
}

func getTodoByID(id string) (*todo, error) {
	for i, t := range todos {
		if t.ID == id {
			return &todos[i], nil
		}
	}
	return nil, errors.New("todo not found")
}
func toggleTodoStatus(c *gin.Context) {
	id := c.Param("id")
	t, err := getTodoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	t.Completed = !t.Completed
	c.IndentedJSON(http.StatusOK, t)
}

func deleteTodoByID(id string) error {
	for i, t := range todos {
		if t.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return errors.New("todo not found")
}

func deleteTodo(c *gin.Context) {
	id := c.Param("id")
	err := deleteTodoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "todo deleted"})

}

func getTodo(c *gin.Context) {
	id := c.Param("id")
	t, err := getTodoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, t)
}

func main() {
	router := gin.Default()
	router.GET("/todos", getTodos)
	router.GET("/todos/:id", getTodo)
	router.PATCH("/todos/:id", toggleTodoStatus)
	router.POST("/todos", addTodo)
	router.DELETE("/todos/:id", deleteTodo)
	router.Run("0.0.0.0:9090")
}
