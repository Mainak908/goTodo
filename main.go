package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/mainak908/simpleTodo/config"
	handlers "github.com/mainak908/simpleTodo/handler"
	"github.com/mainak908/simpleTodo/middleware"
)

func main() {
	client := config.InitDB()
	defer client.Close()

	app := fiber.New()

	app.Post("/register", handlers.Register(client))
	app.Post("/login", handlers.Login(client))

	todo := app.Group("/todos", middleware.Protected())
	todo.Get("/", handlers.GetTodos(client))
	todo.Post("/", handlers.CreateTodo(client))
	todo.Put("/:id", handlers.UpdateTodo(client))
	todo.Delete("/:id", handlers.DeleteTodo(client))

	log.Fatal(app.Listen(":3000"))
}
