package registration

import (
	"chatterbox/accounts"
	"chatterbox/messaging"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// Register registers a user
func Register(rdb *redis.Client) func(*fiber.Ctx) error {
	return websocket.New(func(conn *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		// log.Println(c.Locals("allowed"))  // true
		// log.Println(c.Params("id")) // 123
		// log.Println(c.Query("v"))         // 1.0
		// log.Println(c.Cookies("session")) // ""

		userID := conn.Params("id")
		// channels := c.Query("c")

		user := &accounts.User{
			Id:             userID,
			Name:           userID,
			DisconnectChan: make(chan struct{}),
			MessageChan:    make(chan redis.Message),
		}

		user.Connect(rdb, []string{userID})

		conn.SetCloseHandler(func(code int, text string) error {
			err := user.Disconnect(rdb)
			if err != nil {
				return nil
			}
			return nil
		})

		go messaging.HandleOutbound(*user, conn)
		messaging.HandleInbound(conn, rdb)

	})

}
