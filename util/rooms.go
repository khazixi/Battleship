package util

import (
	"errors"
	"math/rand"
	"net"

	"github.com/khazixi/Battelship/game"
)

type RoomState int

const (
	IDLE RoomState = iota
	OPEN
	FULL
)

type Room struct {
	State       RoomState
	Host        net.Conn
	Participant net.Conn
	Game        *game.Game
}

type RoomList struct {
	M map[int]Room
}

func MakeRoomList() RoomList {
	return RoomList{M: make(map[int]Room)}
}

func (r *RoomList) CreateRoom(host net.Conn) int {
	roomID := rand.Int()
	_, ok := r.M[roomID]

	for ok {
		roomID = rand.Int()
		_, ok = r.M[roomID]
	}

	r.M[roomID] = Room{State: OPEN, Host: host}
	return roomID
}

func (r *RoomList) JoinRoom(roomID int, participant net.Conn) (net.Conn, error) {
	room, ok := r.M[roomID]
	if !ok {
		return nil, errors.New("Attempted to join an empty room")
	}

	room.State = FULL
	room.Participant = participant

	r.M[roomID] = room
	return room.Host, nil
}

func (r *RoomList) GetRooms() []int {
	iterated := 0
	rooms := make([]int, 10)
	for k, v := range r.M {
		if v.State == OPEN {
			rooms = append(rooms, k)
			iterated += 1
		}

		if iterated == 10 {
			break
		}
	}

	return rooms
}

func (r *RoomList) RemoveRoom(roomID int) {
	delete(r.M, roomID)
}

func (r *RoomList) ClearRooms() {
	clear(r.M)
}
