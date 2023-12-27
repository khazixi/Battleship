package util

import (
	"math/rand"
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
	State RoomState
}

func CreateRoom(roomList * sync.Map) uint64 {
	roomID := rand.Uint64()
  _, ok := roomList.Load(roomID)
  
  for ok {
    roomID = rand.Uint64()
    _, ok = roomList.Load(roomID)
  }

  roomList.Store(roomID, Room{State: OPEN})

  return roomID
}

func JoinRoom(roomList * sync.Map, roomID uint64) {
  v, ok := roomList.Load(roomID)
  if !ok {
    return
  }
  room, ok := v.(Room)
  if !ok {
    return
  }
  room.State = FULL
  roomList.Store(roomID, room)
}
