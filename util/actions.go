package util

import (
	"encoding/gob"
	"log"
)

type ActionType int

type Action interface {
	getEvent() ActionType
}

const (
	List ActionType = iota
	Create
	Join
	Delete
)

type JoinAction struct{ RoomID int }

type DeleteAction struct{ RoomID int }

type CreateAction struct{}

type ListAction struct{}

func (a JoinAction) getEvent() ActionType {
	return Join
}

func (c CreateAction) getEvent() ActionType {
	return Create
}

func (l ListAction) getEvent() ActionType {
	return List
}

func (d DeleteAction) getEvent() ActionType {
	return Delete
}

func ActionEncoder(enc *gob.Encoder, a Action) {
	err := enc.Encode(&a)
	if err != nil {
		log.Fatal("Failed to encode", err)
	}
}

func ActionDecoder(dec *gob.Decoder) (Action, error) {
	var a Action
	err := dec.Decode(&a)
	return a, err
}
