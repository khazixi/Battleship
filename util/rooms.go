package util

import (
	"math/rand"
	"net"
	"sync"
)

type RoomState int

const (
	IDLE RoomState = iota
	OPEN
	FULL
)

type RoomList sync.Map

type Room struct {
	State       RoomState
	Host        net.Conn
	Participant net.Conn
}

func CreateRoom(roomList *sync.Map, host net.Conn) int {
	roomID := rand.Int()
	_, ok := roomList.Load(roomID)

	for ok {
		roomID = rand.Int()
		_, ok = roomList.Load(roomID)
	}

	roomList.Store(roomID, Room{State: OPEN, Host: host})

	return roomID
}

func JoinRoom(roomList *sync.Map, roomID int, participant net.Conn) {
	v, ok := roomList.Load(roomID)
	if !ok {
		return
	}
	room, ok := v.(Room)
	if !ok {
		return
	}

	room.State = FULL
  room.Participant = participant

	roomList.Store(roomID, room)
}
