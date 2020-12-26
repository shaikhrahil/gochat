package messaging

import (
	"chatterbox/accounts"
	"log"

	"github.com/go-redis/redis"
	"github.com/gofiber/websocket/v2"
)

const (
	codeAck = iota
	codeMsg
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
func HandleInbound(conn *websocket.Conn, rdb *redis.Client, closeCh chan struct{}) {
loop:
	for {
		select {
		case <-closeCh:
			break loop
		default:
			{
				var parsedMsg Message
				err := conn.ReadJSON(&parsedMsg)
				if err != nil {
					log.Println("Invalid message", err)
					if err := conn.WriteJSON("{test : 'sddslk'}"); err != nil {
						log.Println("Error while error ! : ", err)
					}
				}

				if err == nil {
					if err := rdb.LPush(parsedMsg.Channel, parsedMsg.Message).Err(); err != nil {
						log.Println("Message publish to redis failed !", err.Error())
					}
					if err := rdb.Publish(parsedMsg.Channel, parsedMsg.Message).Err(); err != nil {
						log.Println("Message publish to redis failed !", err.Error())
					}
				}
			}
		}
	}
}
