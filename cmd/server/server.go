package main

import (
	"encoding/gob"
	"io"
	"log"
	"net"
	"sync"

	"github.com/khazixi/Battelship/util"
)

func main() {
	serv, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Printf("An Error occured %v", err)
	}

	defer serv.Close()

	log.Println("Connection Established")

	gob.Register(util.RoomMessage{})
	gob.Register(util.ConfirmationMessage{})
	gob.Register(util.JoinAction{})
	gob.Register(util.ListAction{})
	gob.Register(util.CreateAction{})

	// NOTE: Naive Implimentation prone to race conditions
	msgch := make(chan util.Message)

	go gameLoop(msgch)

	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Fatal("An Error occured ", err)
		}

		go handleConnection(conn, msgch)
	}
}

func handleConnection(conn net.Conn, msgch chan util.Message) {
	log.Println("Launched")
	var roomID int = -1

	defer func() {
    if roomID != -1 {
      msgch <- util.DeleteMessage{RoomID: roomID}
    }
    conn.Close()
  }()

	for {
		decoder := gob.NewDecoder(conn)
		abcd, err := actionDecoder(decoder)

		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		switch currentAction := abcd.(type) {
		case util.CreateAction:
			msgch <- util.CreateMessage{Conn: conn}
			log.Println("This is a creation")
		case util.JoinAction:
			roomID = currentAction.RoomID
			msgch <- util.JoinMessage{RoomID: roomID}
			log.Println("This is an Join")
		case util.ListAction:
			msgch <- util.ListMessage{Conn: conn}
			log.Println("This is a List")
		case util.DeleteAction:
			msgch <- util.DeleteMessage{
				Conn:   conn,
				RoomID: currentAction.RoomID,
			}
		default:
			log.Println("This type is unknown", abcd)
		}
	}
}

func gameLoop(msgch chan util.Message) {
	roomList := sync.Map{}
	for {
		select {
		case m := <-msgch:
			switch message := m.(type) {
			case util.CreateMessage:
				util.CreateRoom(&roomList, message.Conn)
			case util.JoinMessage:
				host, err := util.JoinRoom(&roomList, message.RoomID, message.Conn)
				enc := gob.NewEncoder(message.Conn)
				if err != nil {
					messageEncoder(enc, util.ConfirmationMessage{
						Joined: false,
						RoomID: message.RoomID,
					})
				}
				host_enc := gob.NewEncoder(host)
				messageEncoder(enc, util.ConfirmationMessage{
					Joined: true,
					RoomID: message.RoomID,
				})
				messageEncoder(host_enc, util.ConfirmationMessage{
					Joined: true,
					RoomID: message.RoomID,
				})
			case util.DeleteMessage:
				roomList.Delete(message.RoomID)
			case util.ListMessage:
				rooms := util.GetRooms(&roomList)
				enc := gob.NewEncoder(message.Conn)
				messageEncoder(enc, util.RoomsMessage{Rooms: rooms})
			case util.ClearMessage:
				roomList.Range(func(key, value any) bool {
					roomList.Delete(key)
					return true
				})
			}

		}
	}
}

func messageEncoder(enc *gob.Encoder, m util.Message) {
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
