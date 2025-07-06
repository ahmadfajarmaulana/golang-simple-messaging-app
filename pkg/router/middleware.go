package router

import (
	"log"
	"simple-messaging-app/app/repository"
	"simple-messaging-app/pkg/jwt_token"
	"simple-messaging-app/pkg/response"
	"time"

	"github.com/gofiber/fiber/v2"
)

func MiddlewareValidateAuth(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		log.Println("authorization is empty")
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	_, err := repository.GetUserSessionByToken(ctx.Context(), auth)
	if err != nil {
		log.Println("failed to get user session by token on db: ", err)
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwt_token.ValidateToken(ctx.Context(), auth)
	if err != nil {
		log.Println(err)
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token is expired")
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.FullName)

	return ctx.Next()
}

func MiddlewareRefreshToken(ctx *fiber.Ctx) error {
	auth := ctx.Get("authorization")
	if auth == "" {
		log.Println("authorization is empty")
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	claim, err := jwt_token.ValidateToken(ctx.Context(), auth)
	if err != nil {
		log.Println(err)
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	if time.Now().Unix() > claim.ExpiresAt.Unix() {
		log.Println("jwt token is expired")
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, "unauthorized", nil)
	}

	ctx.Locals("username", claim.Username)
	ctx.Locals("full_name", claim.FullName)
	return ctx.Next()
}
