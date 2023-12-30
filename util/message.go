package util

import "net"

type MessageType int

const (
	ROOM MessageType = iota
	CONFIRM
	CREATE
	JOIN
	CLOSED
  LIST
)

type Message interface {
	getType() MessageType
}

type RoomMessage struct {
	RoomID int
}

type JoinMessage struct {
	RoomID int
	Conn   net.Conn
}

type CloseMessage struct {
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

func (r RoomMessage) getType() MessageType {
	return ROOM
}

func (c ConfirmationMessage) getType() MessageType {
	return CONFIRM
}

func (c CloseMessage) getType() MessageType {
	return CLOSED
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
