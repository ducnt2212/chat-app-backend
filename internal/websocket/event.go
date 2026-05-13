package websocket

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventTyping      = "typing"
)

type Event struct {
	Type    string `json:"type"`
	RoomID  int    `json:"room_id"`
	Payload any    `json:"payload"`
}
