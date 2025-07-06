package ws

import (
	"context"
	"fmt"
	"log"
	"simple-messaging-app/app/models"
	"simple-messaging-app/app/repository"
	"simple-messaging-app/pkg/env"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func ServeWSMessaging(app *fiber.App) {
	var clients = make(map[*websocket.Conn]bool)
	var broadcast = make(chan models.MessagePayload)

	app.Get("/message/v1/send", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()

		clients[c] = true
		for {
			var msg models.MessagePayload
			if c.ReadJSON(&msg) != nil {
				log.Println("error payload: ", msg)
				break
			}

			msg.Date = time.Now()
			err := repository.InsertNewMessage(context.Background(), &msg)
			if err != nil {
				log.Println(err)
				break
			}
			broadcast <- msg
		}
	}))

	go func() {
		for {
			msg := <-broadcast
			for client := range clients {
				if err := client.WriteJSON(msg); err != nil {
					log.Println("Failed to write json: ", msg)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_SOCKET_PORT", "8080"))))
}
