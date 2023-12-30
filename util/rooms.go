package util

import (
	"errors"
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
	Game        Game
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

func JoinRoom(roomList *sync.Map, roomID int, participant net.Conn) (net.Conn, error) {
	v, ok := roomList.Load(roomID)
	if !ok {
		return nil, errors.New("Attempted to join an empty room")
	}
	room, ok := v.(Room)
	if !ok {
		return nil, errors.New("Room contained an incorrect type")
	}

	room.State = FULL
	room.Participant = participant

	roomList.Store(roomID, room)
	return room.Host, nil
}

func GetRooms(roomList *sync.Map) []int {
	iterated := 0
	rooms := make([]int, 10)
	roomList.Range(func(key, value any) bool {
		pk, ok := key.(int)
		if !ok {
			return true
		}

		_, ok = value.(Room)
		if !ok {
			return true
		}

		rooms = append(rooms, pk)

		return iterated < 11
	})

	return rooms
}
