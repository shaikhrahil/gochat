package messaging

import (
	"encoding/json"
	"log"

	"github.com/go-redis/redis"
	"github.com/gofiber/websocket/v2"
)

// HandleOutbound handles outgoing messages from server
func HandleOutbound(msgChan chan redis.Message, conn *websocket.Conn, formatMessage func(msg string) Message) error {
	for m := range msgChan {
		if err := conn.WriteJSON(formatMessage(m.Payload)); err != nil {
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
					if err := conn.WriteJSON("{ From : 'Sys', Message : 'Invalid Message'}"); err != nil {
						log.Println("Error while error ! : ", err)
					}
				}

				msg, err := json.Marshal(parsedMsg)

				if err != nil {
					errMsg, _ := json.Marshal(Message{From: From{ID: "0", Name: "Sys"}, Message: "Invalid Message"})
					if err := conn.WriteJSON(errMsg); err != nil {
						log.Println("Error while error ! : ", err)
					}
				}

				if err := rdb.LPush(parsedMsg.Channel, msg).Err(); err != nil {
					log.Println("Message publish to redis failed !", err.Error())
				}
				if err := rdb.Publish(parsedMsg.Channel, parsedMsg.Message).Err(); err != nil {
					log.Println("Message publish to redis failed !", err.Error())
				}
			}
		}
	}
}
