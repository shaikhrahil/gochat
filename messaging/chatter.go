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

// HandleInbound handles incoming messages
func HandleInbound(user accounts.User, conn *websocket.Conn) error {
	for m := range user.MessageChan {
		if err := conn.WriteJSON(m.Payload); err != nil {
			log.Println("Some error occured while writing to socket", err)
		}
	}
	return nil
}

// HandleOutbound handles outgoing messages
func HandleOutbound(conn *websocket.Conn, rdb *redis.Client) {

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
				log.Println("Connection failed !")
				return
			}

			err := rdb.Publish(parsedMsg.Channel, parsedMsg.Message)
			if someErr := err.Err(); someErr != nil {
				log.Println("Message publish tot redis failed !", parsedMsg)
				conn.WriteJSON(someErr.Error())
			}
		}
	}
}
