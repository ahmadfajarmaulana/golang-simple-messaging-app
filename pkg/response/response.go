package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	ResponseSuccess struct {
		Status  bool        `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	ResponseError struct {
		Status  bool        `json:"status"`
		Message string      `json:"message"`
		Errors  interface{} `json:"errors"`
	}
)

func SendSuccessResponse(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(ResponseSuccess{
		Status:  true,
		Message: "success",
		Data:    data,
	})
}

func SendErrorResponse(ctx *fiber.Ctx, httpCode int, message string, err interface{}) error {
	return ctx.Status(httpCode).JSON(ResponseError{
		Status:  false,
		Message: message,
		Errors:  err,
	})
}

func SendValidationResponse(ctx *fiber.Ctx, err interface{}) error {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			if jsonTag := fieldErr.StructField(); jsonTag != "" {
				fieldName = fieldErr.Field()
			}
			errors[fieldName] = fmt.Sprintf("field is %s", fieldErr.Tag())
		}
	} else if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(ResponseError{
			Status:  false,
			Message: "Validation failed",
			Errors:  err,
		})
	} else {
		return ctx.Status(fiber.StatusInternalServerError).JSON(ResponseError{
			Status:  false,
			Message: "Internal server error",
			Errors:  err,
		})
	}

	return ctx.Status(fiber.StatusBadRequest).JSON(ResponseError{
		Status:  false,
		Message: "Validation failed",
		Errors:  errors,
	})
}
