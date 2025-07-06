package controllers

import (
	"log"
	"simple-messaging-app/app/repository"
	"simple-messaging-app/pkg/response"

	"github.com/gofiber/fiber/v2"
)

func GetHistory(ctx *fiber.Ctx) error {
	resp, err := repository.GetAllMessage(ctx.Context())
	if err != nil {
		log.Println(err)
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}
	return response.SendSuccessResponse(ctx, resp)
}
