package main

import (
	"encoding/gob"
	"io"
	"log"
	"net"
	"sync"

	"github.com/khazixi/Battelship/util"
)


type Action struct {
	A, B, C, D byte
}

type Inaction struct {
	E int
}

var roomList sync.Map = sync.Map{}

func main() {
	serv, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("An Error occured %s", err.Error)
	}

	defer serv.Close()

	log.Println("Connection Established")

  gob.Register(util.RoomMessage{})
	gob.Register(util.JoinAction{})
	gob.Register(util.ListAction{})
	gob.Register(util.CreateAction{})

  // NOTE: Naive Implimentation prone to race conditions

	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Fatal("An Error occured ", err)
		}

    handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
    decoder := gob.NewDecoder(conn)
    encoder := gob.NewEncoder(conn)
		abcd, err := actionDecoder(decoder)

		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
      return
		}

    switch currentAction := abcd.(type) {
		case util.CreateAction:
      log.Println("This is a creation")
      roomID := util.CreateRoom(&roomList)
      // err = encoder.Encode(util.CreateRoomMessage(roomID))
      messageEncoder(encoder, util.CreateRoomMessage(roomID))
      if err != nil {
        log.Println("Error with sending the message")
      }
		case util.JoinAction:
      util.JoinRoom(&roomList, currentAction.RoomID)
			log.Println("This is an Join")
		case util.ListAction:
			log.Println("This is a List")
		default:
			log.Println("This type is unknown", abcd)
		}
	}
}

func messageEncoder(enc * gob.Encoder, m util.Message) {
  err := enc.Encode(&m)
  if err != nil {
    log.Fatal("Failed to encode", err)
  }
}

func actionDecoder(dec *gob.Decoder) (util.Action, error) {
	var a util.Action
	err := dec.Decode(&a)
	return a, err
}
