package util

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/khazixi/Battelship/game"
)

type msgType uint8

// These constants are used to determine the type of message or action needed to be communicated
const (
	status msgType = iota
	join
	create
	list
	clear
	delete
	initialize
	leave
	exit
)

// INFO: This is used for sending

type ServerMsg interface {
  serverMsg()
}

// Communicates the result of the client sent message back to the client.
// This should be used only on the client.
type StatusMsg struct {
	msgType
	action msgType
	status bool
}

// This message is sent back for cases of list actions.
type RoomMsg struct {
	msgType
	rooms []int
}

// No-Op boilerplate for interfaces because I have no idea what else to do
func (m StatusMsg) serverMsg() {}

func (r RoomMsg) serverMsg() {}

func ServerMsgEncoder(enc * gob.Encoder, s ServerMsg) {
  err := enc.Encode(&s)
  if err!= nil {
    log.Fatal("Failed to Encode: ", err)
  }
}

func ServerMsgDecoder(dec * gob.Decoder) (ServerMsg, error) {
  var msg ServerMsg
  err := dec.Decode(&msg)
  return msg, err
}

// INFO: This struct is for forwarding messages from connection to game server

// Server message used for when rooms aren't important
type InternalMsg struct {
	msgType
	room int
	conn net.Conn
}

// INFO: Below are the structs that the client sends to the server

type ClientMsg interface {
  clientMsg()
}

// Communicates from the client to the server.
type ActionMsg struct {
	msgType
	action msgType

	// If the value of room is negative then ignore the value of room
	room int
}

type InitMsg struct {
	msgType
	conn     net.Conn
	transmit [5]game.Transmit
}

func ClientMsgEncoder(enc * gob.Encoder, s ClientMsg) {
  err := enc.Encode(&s)
  if err!= nil {
    log.Fatal("Failed to Encode: ", err)
  }
}

func ClientMsgDecoder(dec * gob.Decoder) (ClientMsg, error) {
  var msg ClientMsg
  err := dec.Decode(&msg)
  return msg, err
}

func (m ActionMsg) clientMsg() {}

func (m InitMsg) clientMsg() {}
