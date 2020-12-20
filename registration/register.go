package registration

import (
	"chatterbox/accounts"
	"chatterbox/messaging"
	"log"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Register registers a user
func Register(rdb *redis.Client) func(*fiber.Ctx) error {
	return websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		// log.Println(c.Locals("allowed"))  // true
		// log.Println(c.Params("id")) // 123
		// log.Println(c.Query("v"))         // 1.0
		// log.Println(c.Cookies("session")) // ""

		userID := c.Params("id")
		// channels := c.Query("c")

		user := &accounts.User{
			Id:             userID,
			Name:           userID,
			DisconnectChan: make(chan struct{}),
			MessageChan:    make(chan redis.Message),
		}

		user.Connect(rdb, userID)

		c.SetCloseHandler(func(code int, text string) error {
			err := user.Disconnect(rdb)
			if err != nil {
				return nil
			}
			return nil
		})

		go messaging.HandleInbound(*user, c)

		var (
			mt  int
			msg []byte
			err error
		)

		for {
			mt, msg, err = c.ReadMessage()
			if err != nil {
				log.Println(mt, msg, err)
			}
			if err == nil {
				var parsedMsg messaging.Message

				if err := c.ReadJSON(&parsedMsg); err != nil {
					_ = c.WriteJSON("msg{Err: err.Error()}")
					return
				}

				err := rdb.Publish(parsedMsg.Channel, parsedMsg.Message)
				if someErr := err.Err(); someErr != nil {
					c.WriteJSON(someErr.Error())
				}
			}
		}

	})

}
