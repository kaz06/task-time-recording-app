package main

import (
	"log"
	"net/http"

	"go-backend/internal/db"
	"go-backend/internal/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
			http.MethodOptions},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Handlers
	c, err := handlers.NewContainer(database)
	if err != nil {
		log.Fatalf("Failed to create handler container: %v", err)
	}
	e.Use(c.AuthMiddleware)

	// Routes
	e.POST("/v1/users", c.UsersPost)
	e.GET("/v1/users", c.UsersGet)
	e.GET("/v1/tasks", c.TasksGet)
	e.GET("/v1/tasktimebytag", c.TasksTagGet)
	e.POST("/v1/tasks", c.TasksPost)
	e.DELETE("/v1/tasks/:taskId", c.TasksTaskIdDelete)
	e.GET("/v1/tasks/:taskId", c.TasksTaskIdGet)
	e.PUT("/v1/tasks/:taskId", c.TasksTaskIdPut)

	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	log.Fatal(http.ListenAndServe(":8080", e))
}
