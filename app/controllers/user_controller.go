package controllers

import (
	"fmt"

	"simple-messaging-app/app/models"
	"simple-messaging-app/app/repository"
	"simple-messaging-app/pkg/jwt_token"
	"simple-messaging-app/pkg/response"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func Register(ctx *fiber.Ctx) error {
	user := &models.User{}
	err := ctx.BodyParser(user)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request: %v", err)
		fmt.Println(errResponse)
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	err = user.Validate()
	if err != nil {
		fmt.Println("Failed to validate request: ", err)
		return response.SendValidationResponse(ctx, err)
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errResponse := fmt.Errorf("failed to endcrypt the password: %v", err)
		fmt.Println(errResponse)
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	user.Password = string(hashPassword)

	err = repository.InsertNewUser(ctx.Context(), user)
	if err != nil {
		errResponse := fmt.Errorf("failed create data to database: %v", err)
		fmt.Println(errResponse)
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	res := user
	res.Password = ""

	return response.SendSuccessResponse(ctx, res)
}

func Login(ctx *fiber.Ctx) error {
	//parse request
	loginRequest := &models.LoginRequest{}
	err := ctx.BodyParser(loginRequest)
	if err != nil {
		errResponse := fmt.Errorf("failed to parse request: %v", err)
		fmt.Println(errResponse)
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, errResponse.Error(), nil)
	}

	err = loginRequest.Validate()
	if err != nil {
		fmt.Println("Failed to validate request: ", err)
		return response.SendValidationResponse(ctx, err)
	}

	user, err := repository.GetUserByUsername(ctx.Context(), loginRequest.Username)
	if err != nil {
		errResponse := fmt.Errorf("failed to get user by username: %v", err)
		fmt.Println(errResponse)
		return response.SendErrorResponse(ctx, fiber.StatusNotFound, errResponse.Error(), nil)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		errResponse := fmt.Errorf("failed to compare password : %v", err)
		fmt.Println(errResponse)

		passwordError := map[string]string{
			"password": "invalid password",
		}

		return response.SendValidationResponse(ctx, passwordError)
	}

	token, err := jwt_token.GenerateToken(ctx.Context(), user.Username, user.FullName, "token")
	if err != nil {
		errResponse := fmt.Errorf("failed to generate token: %v", err)
		fmt.Println(errResponse)
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	refreshToken, err := jwt_token.GenerateToken(ctx.Context(), user.Username, user.FullName, "refresh_token")
	if err != nil {
		errResponse := fmt.Errorf("failed to generate refresh_token: %v", err)
		fmt.Println(errResponse)
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, "internal server error", nil)
	}

	res := models.LoginResponse{
		Username:     user.Username,
		FullName:     user.FullName,
		Token:        token,
		RefreshToken: refreshToken,
	}
	return response.SendSuccessResponse(ctx, res)
}
