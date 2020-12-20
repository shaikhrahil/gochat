package messaging

import (
	"chatterbox/accounts"
	"fmt"

	"github.com/gofiber/websocket/v2"
)

// HandleInbound handles incoming messages
func HandleInbound(user accounts.User, conn *websocket.Conn) error {
	for m := range user.MessageChan {
		if err := conn.WriteJSON(m.Payload); err == nil {
			fmt.Println(err)
		}
	}
	return nil
}
