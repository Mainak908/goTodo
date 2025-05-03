package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mainak908/simpleTodo/ent"
	"github.com/mainak908/simpleTodo/ent/user"
	"golang.org/x/crypto/bcrypt"
)

type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input AuthInput
		if err := c.BodyParser(&input); err != nil {
			return err
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		_, err := client.User.Create().
			SetEmail(input.Email).
			SetPassword(string(hash)).
			Save(context.Background())
		if err != nil {
			return fiber.ErrConflict
		}
		return c.SendStatus(fiber.StatusCreated)
	}
}

func Login(client *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input AuthInput
		if err := c.BodyParser(&input); err != nil {
			return err
		}
		u, err := client.User.
			Query().
			Where(user.EmailEQ(input.Email)).
			Only(context.Background())
		if err != nil || bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)) != nil {
			return fiber.ErrUnauthorized
		}
		claims := jwt.MapClaims{
			"id":  u.ID,
			"exp": time.Now().Add(time.Hour * 72).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		s, _ := token.SignedString([]byte("secret"))

		
		return c.JSON(fiber.Map{"token": s})

	}
}