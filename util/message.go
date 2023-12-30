package util

import "net"

type MessageType int

const (
	ROOM MessageType = iota
	CONFIRM
	CREATE
	JOIN
  CLEAR
  LIST
  ROOMS
  DELETE
)

type Message interface {
	getType() MessageType
}

type RoomMessage struct {
	RoomID int
}

type ClearMessage struct{}

type JoinMessage struct {
	RoomID int
	Conn   net.Conn
}

type DeleteMessage struct {
	RoomID int
	Conn   net.Conn
}

type ConfirmationMessage struct {
	Joined bool
	RoomID int
}

type CreateMessage struct {
	Conn  net.Conn
}

type ListMessage struct {
  Conn net.Conn
}

type RoomsMessage struct {
  Rooms []int
}

func (r RoomMessage) getType() MessageType {
	return ROOM
}

func (c ConfirmationMessage) getType() MessageType {
	return CONFIRM
}

func (c ClearMessage) getType() MessageType {
	return CLEAR
}

func (c CreateMessage) getType() MessageType {
	return CREATE
}

func (c JoinMessage) getType() MessageType {
	return JOIN
}

func (l ListMessage) getType() MessageType {
  return LIST
}

func (r RoomsMessage) getType() MessageType {
  return ROOMS
}

func (d DeleteMessage) getType() MessageType {
  return DELETE
}
