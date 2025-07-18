package repository

import (
	"context"
	"fmt"
	"simple-messaging-app/app/models"
	"simple-messaging-app/pkg/database"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func InsertNewMessage(ctx context.Context, data *models.MessagePayload) error {
	_, err := database.MongoDB.InsertOne(ctx, data)
	return err
}

func GetAllMessage(ctx context.Context) ([]models.MessagePayload, error) {
	var (
		err  error
		resp []models.MessagePayload
	)

	cursor, err := database.MongoDB.Find(ctx, bson.D{})

	if err != nil {
		return resp, fmt.Errorf("failed to find message: %v", err)
	}

	for cursor.Next(ctx) {
		var payload models.MessagePayload
		err = cursor.Decode(&payload)
		if err != nil {
			return resp, fmt.Errorf("failed to decode payload: %v", err)
		}
		resp = append(resp, payload)
	}
	return resp, nil

}
