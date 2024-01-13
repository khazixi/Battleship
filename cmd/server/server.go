package main

import (
	"encoding/gob"
	"io"
	"log"
	"net"

	"github.com/khazixi/Battelship/game"
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
	gob.Register(util.RoomsMessage{})

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
		log.Println("Disconnected")
		if roomID != -1 {
			msgch <- util.DeleteMessage{RoomID: roomID}
		}
		conn.Close()
	}()

	for {
		decoder := gob.NewDecoder(conn)
		abcd, err := util.ActionDecoder(decoder)

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
		case util.LeaveAction:
			msgch <- util.LeaveMessage{
				Conn:   conn,
				RoomID: currentAction.RoomID,
			}
		default:
			log.Println("This type is unknown", abcd)
		}
	}
}

func gameLoop(msgch chan util.Message) {
	var roomList util.RoomList = util.MakeRoomList()
	for {
		select {
		case m := <-msgch:
			switch message := m.(type) {
			case util.CreateMessage:
				roomList.CreateRoom(message.Conn)
			case util.JoinMessage:
				host, err := roomList.JoinRoom(message.RoomID, message.Conn)
				enc := gob.NewEncoder(message.Conn)
				if err != nil {
					util.MessageEncoder(enc, util.ConfirmationMessage{
						Joined: false,
						RoomID: message.RoomID,
					})
				}
				host_enc := gob.NewEncoder(host)
				util.MessageEncoder(enc, util.ConfirmationMessage{
					Joined: true,
					RoomID: message.RoomID,
				})
				util.MessageEncoder(host_enc, util.ConfirmationMessage{
					Joined: true,
					RoomID: message.RoomID,
				})
				// WARNING: Maybe move this to its own message
				roomList.M[message.RoomID].Game.PlayerTurn = game.PLAYER1
			case util.DeleteMessage:
				roomList.RemoveRoom(message.RoomID)
			case util.ListMessage:
				rooms := roomList.GetRooms()
				enc := gob.NewEncoder(message.Conn)
				util.MessageEncoder(enc, util.RoomsMessage{Rooms: rooms})
			case util.ClearMessage:
				roomList.ClearRooms()
			case util.LeaveMessage:
        host := roomList.M[message.RoomID].Host
        participant := roomList.M[message.RoomID].Participant
        hostEnc := gob.NewEncoder(host)
        partEnc := gob.NewEncoder(participant)
        hostEnc.Encode(util.ExitedMessage{})
        partEnc.Encode(util.ExitedMessage{})
        delete(roomList.M, message.RoomID)
			case util.InitializerMessage:
				for _, v := range message.Transmit {
					if message.Conn == roomList.M[message.Room].Host {
						roomList.M[message.Room].Game.P1.Place(v.Coordinate, v.Piece, v.Direction)
					} else if message.Conn == roomList.M[message.Room].Participant {
						roomList.M[message.Room].Game.P2.Place(v.Coordinate, v.Piece, v.Direction)
					}
				}
			}

		}
	}
}
