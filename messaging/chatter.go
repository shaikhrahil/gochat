package messaging

import (
	"chatterbox/accounts"
	"log"

	"github.com/go-redis/redis"
	"github.com/gofiber/websocket/v2"
)

const (
	connectCode = iota
	disconnectCode
	messageCode
)

// HandleOutbound handles outgoing messages from server
func HandleOutbound(user accounts.User, conn *websocket.Conn) error {
	for m := range user.MessageChan {
		if err := conn.WriteJSON(m.Payload); err != nil {
			log.Println("Some error occured while writing to socket", err)
		}
	}
	return nil
}

// HandleInbound handles incoming messages to the server
func HandleInbound(conn *websocket.Conn, rdb *redis.Client) {

	var (
		mt  int
		msg []byte
		err error
	)

	mt, msg, err = conn.ReadMessage()
	if err != nil {
		log.Println(mt, msg, err)
	}
	if err == nil {
		for {
			var parsedMsg Message
			if err := conn.ReadJSON(&parsedMsg); err != nil {
				log.Println("Invalid message", err)
				if err := conn.WriteJSON("{test : 'sddslk'}"); err != nil {
					log.Println("Error while error ! : ", err)
				}
			}

			if err := rdb.Publish(parsedMsg.Channel, parsedMsg.Message).Err(); err != nil {
				log.Println("Message publish tot redis failed !", err.Error())
			}
		}
	}
}
