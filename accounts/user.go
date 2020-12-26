package accounts

import (
	"log"

	"github.com/go-redis/redis"
)

// User model
type User struct {
	Id             string
	Name           string
	DisconnectChan chan struct{}
	MessageChan    chan redis.Message
	ChannelHandler *redis.PubSub
}

func (u *User) Disconnect(rdb *redis.Client) error {
	err := u.ChannelHandler.Unsubscribe()
	if err != nil {
		return err
	}
	log.Println("Unsubscribed " + u.Id + " from redis")
	if userExists := rdb.Exists(u.Id); userExists != nil {
		rdb.Del(u.Id)
	}
	close(u.DisconnectChan)
	log.Println("Removed " + u.Id + " from redis")
	return nil
}

func (u *User) Connect(rdb *redis.Client, channels []string) error {
	log.Println("Subscribed to ", channels)
	pubSub := rdb.Subscribe(channels...)
	for _, ch := range channels {
		test := rdb.LRange(ch, 0, -1)
		if test.Err() == nil {
			for _, msg := range test.Val() {
				rdb.Publish(ch, msg)
			}
		}
	}
	u.ChannelHandler = pubSub

	go func() {
		for {
			select {
			case msg, ok := <-pubSub.Channel():
				log.Println("Received message via pubsub at ", channels)
				if !ok {
					log.Println("Message not ok for ", channels, msg)
					return
				}
				log.Println("Sent message to ", channels)
				rdb.LPush(msg.Channel, msg.Payload)
				u.MessageChan <- *msg
			case <-u.DisconnectChan:
				log.Println("Disconnected channel triggered fro ", channels)
				return
			}
		}
	}()
	return nil
}
