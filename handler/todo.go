package handler

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mainak908/simpleTodo/ent"
	"github.com/mainak908/simpleTodo/ent/todo"
	"github.com/mainak908/simpleTodo/ent/user"
)

type TodoInput struct {
	Title string `json:"title"`
}

func getUserID(c *fiber.Ctx) (int, error) {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := int(claims["id"].(float64))
	return id, nil
}

func GetTodos(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := getUserID(c)
		todos, _ := client.Todo.Query().Where(todo.HasUserWith(user.ID(uid))).All(context.Background())
		return c.JSON(todos)
	}
}

func CreateTodo(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uid, _ := getUserID(c)
		var input TodoInput
		if err := c.BodyParser(&input); err != nil {
			return err
		}
		t, _ := client.Todo.Create().SetTitle(input.Title).SetUserID(uid).Save(context.Background())
		return c.JSON(t)
	}
}

func UpdateTodo(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		var input TodoInput
		if err := c.BodyParser(&input); err != nil {
			return err
		}
		t, err := client.Todo.UpdateOneID(id).SetTitle(input.Title).Save(context.Background())
		if err != nil {
			return fiber.ErrNotFound
		}
		return c.JSON(t)
	}
}

func DeleteTodo(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		if err := client.Todo.DeleteOneID(id).Exec(context.Background()); err != nil {
			return fiber.ErrNotFound
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}