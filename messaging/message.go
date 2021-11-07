package messaging

// Message model
type Message struct {
	Channel string
	Message string
	From    From
}

type From struct {
	ID   string
	Name string
}
