package repository

import (
	"context"

	"simple-messaging-app/app/models"
	"simple-messaging-app/pkg/database"
)

func InsertNewUser(ctx context.Context, user *models.User) error {
	return database.DB.Create(user).Error
}

func GetUserByUsername(ctx context.Context, username string) (models.User, error) {
	var (
		res models.User
		err error
	)
	err = database.DB.Where("username = ?", username).Last(&res).Error
	return res, err
}
