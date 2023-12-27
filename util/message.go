package util

type MessageType int

const (
	ROOM MessageType = iota
	CONFIRM
)

type Message interface {
	getType() MessageType
}

type RoomMessage struct {
	messageType MessageType
	RoomID      int
}

type ConfirmationMessage struct {
	messageType MessageType
	Joined      bool
	RoomID      int
}

func CreateRoomMessage(id int) RoomMessage {
	return RoomMessage{
		messageType: ROOM,
		RoomID:      id,
	}
}

func CreateConfirmationMessage(confirmed bool, id int) ConfirmationMessage {
	return ConfirmationMessage{
		messageType: CONFIRM,
		Joined:      confirmed,
		RoomID:      id,
	}
}

func (r RoomMessage) getType() MessageType {
	return r.messageType
}

func (c ConfirmationMessage) getType() MessageType {
	return c.messageType
}
