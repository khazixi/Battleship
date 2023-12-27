package util

type MessageType int

const (
  ROOM MessageType = iota
)

type Message interface {
	getType() MessageType
}

type RoomMessage struct {
	messageType MessageType
	RoomID      uint64
}

func CreateRoomMessage(id uint64) RoomMessage {
  return RoomMessage{
    messageType: ROOM,
    RoomID: id,
  }
}

func (r RoomMessage) getType() MessageType {
  return r.messageType
}
