package util

import (
	"encoding/gob"
	"log"
	"net"

	"github.com/khazixi/Battelship/game"
)

type MsgType uint8

// These constants are used to determine the type of message or action needed to be communicated
const (
	Status MsgType = iota
	Join
	Create
	List
	Clear
	Delete
	Initialize
	Leave
	Exit
)

// INFO: This is used for sending

type ServerMsg interface {
	serverMsg()
}

// Communicates the result of the client sent message back to the client.
// This should be used only on the client.
type StatusMsg struct {
	MsgType
	Action MsgType
	Status bool
	Room   int
}

// This message is sent back for cases of list actions.
type RoomMsg struct {
	MsgType
	Rooms []int
}

// No-Op boilerplate for interfaces because I have no idea what else to do
func (m StatusMsg) serverMsg() {}

func (r RoomMsg) serverMsg() {}

func ServerMsgEncoder(enc *gob.Encoder, s ServerMsg) {
	err := enc.Encode(&s)
	if err != nil {
		log.Fatal("Failed to Encode: ", err)
	}
}

func ServerMsgDecoder(dec *gob.Decoder) (ServerMsg, error) {
	var msg ServerMsg
	err := dec.Decode(&msg)
	return msg, err
}

// INFO: This struct is for forwarding messages from connection to game server

// Server message used for when rooms aren't important

type InternalMsg interface {
	internalMsg()
}

type ConnectionMsg struct {
	MsgType
	Room int
	Conn net.Conn
}

type StartMsg struct {
	MsgType
	Conn     net.Conn
	Room     int
	Transmit [5]game.Transmit
}

func (c ConnectionMsg) internalMsg() {}
func (c StartMsg) internalMsg()      {}

// INFO: Below are the structs that the client sends to the server

type ClientMsg interface {
	clientMsg()
}

// Communicates from the client to the server.
type ActionMsg struct {
	MsgType
	Action MsgType

	// If the value of room is negative then ignore the value of room
	Room int
}

type InitMsg struct {
	MsgType
	Room     int
	Transmit [5]game.Transmit
}

func ClientMsgEncoder(enc *gob.Encoder, s ClientMsg) {
	err := enc.Encode(&s)
	if err != nil {
		log.Fatal("Failed to Encode: ", err)
	}
}

func ClientMsgDecoder(dec *gob.Decoder) (ClientMsg, error) {
	var msg ClientMsg
	err := dec.Decode(&msg)
	return msg, err
}

func (m ActionMsg) clientMsg() {}

func (m InitMsg) clientMsg() {}
