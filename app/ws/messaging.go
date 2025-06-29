package ws

import (
	"fmt"
	"log"
	"simple-messaging-app/pkg/env"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type MessagePayload struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func ServeWSMessaging(app *fiber.App) {
	var clients = make(map[*websocket.Conn]bool)
	var broadcast = make(chan MessagePayload)

	app.Get("/message/v1/send", websocket.New(func(c *websocket.Conn) {
		defer func() {
			c.Close()
			delete(clients, c)
		}()

		clients[c] = true
		for {
			var msg MessagePayload
			if c.ReadJSON(&msg) != nil {
				fmt.Println("error payload: ", msg)
			}

			broadcast <- msg
		}
	}))

	go func() {
		for {
			msg := <-broadcast
			for client := range clients {
				if err := client.WriteJSON(msg); err != nil {
					fmt.Println("Failed to write json: ", msg)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}()

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", env.GetEnv("APP_HOST", "localhost"), env.GetEnv("APP_SOCKET_PORT", "8080"))))
}
