package messaging

// Message model
type Message struct {
	Code    int
	Channel string
	Message string
}

type SubscriptionMessage struct {
	Channels []string
}
