package ws

const (
	BlocksChannel       = "blocks"
	TransactionsChannel = "transactions"

	SubscribeMsgType   = "subscribe"
	UnsubscribeMsgType = "unsubscribe"
)

type (
	msg struct {
		Type    string `json:"type"`
		Channel string `json:"channel"`
	}

	Broadcast struct {
		Channel string      `json:"channel"`
		Data    interface{} `json:"data"`
	}
)
