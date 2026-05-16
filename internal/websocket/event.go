package websocket

type Type string

const (
	EventSendMessage Type = "send_message"
	EventNewMessage  Type = "new_message"
	EventIsTyping    Type = "is_typing"
	EventIsOnline    Type = "is_online"
	EventIsOffline   Type = "is_offline"
)

type Event struct {
	Type    Type `json:"type"`
	RoomID  int  `json:"room_id"`
	Payload any  `json:"payload"`
}
