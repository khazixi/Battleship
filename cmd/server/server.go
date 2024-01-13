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

	gob.Register(util.StatusMsg{})
	gob.Register(util.RoomMsg{})
	gob.Register(util.ActionMsg{})
	gob.Register(util.InitMsg{})

	// NOTE: Naive Implimentation prone to race conditions
	msgch := make(chan util.InternalMsg)

	go gameLoop(msgch)

	for {
		conn, err := serv.Accept()
		if err != nil {
			log.Fatal("An Error occured ", err)
		}

		go handleConnection(conn, msgch)
	}
}

func handleConnection(conn net.Conn, msgch chan util.InternalMsg) {
	log.Println("Launched")
	var roomID int = -1

	defer func() {
		log.Println("Disconnected")
		if roomID != -1 {
			msgch <- util.ConnectionMsg{
				MsgType: util.Delete,
        Room: roomID,
			}
		}
		conn.Close()
	}()

	for {
		decoder := gob.NewDecoder(conn)
		abcd, err := util.ClientMsgDecoder(decoder)

		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			return
		}

		switch currentAction := abcd.(type) {
		case util.ActionMsg:
			msgch <- util.ConnectionMsg{
				MsgType: currentAction.Action,
				Room:    currentAction.Room,
				Conn:    conn,
			}

		case util.InitMsg:
			// TODO: Impliment game features
		}

		// switch currentAction := abcd.(type) {
		// case util.CreateAction:
		// 	msgch <- util.CreateMessage{Conn: conn}
		// 	log.Println("This is a creation")
		// case util.JoinAction:
		// 	roomID = currentAction.RoomID
		// 	msgch <- util.JoinMessage{RoomID: roomID}
		// 	log.Println("This is an Join")
		// case util.ListAction:
		// 	msgch <- util.ListMessage{Conn: conn}
		// 	log.Println("This is a List")
		// case util.DeleteAction:
		// 	msgch <- util.DeleteMessage{
		// 		Conn:   conn,
		// 		RoomID: currentAction.RoomID,
		// 	}
		// case util.LeaveAction:
		// 	msgch <- util.LeaveMessage{
		// 		Conn:   conn,
		// 		RoomID: currentAction.RoomID,
		// 	}
		// default:
		// 	log.Println("This type is unknown", abcd)
		// }
	}
}

func gameLoop(msgch chan util.InternalMsg) {
	var roomList util.RoomList = util.MakeRoomList()
	for {
		select {
		case message := <-msgch:
			switch message := message.(type) {
			case util.StartMsg:
				// TODO: Impliment
				for _, v := range message.Transmit {
					if message.Conn == roomList.M[message.Room].Host {
						roomList.M[message.Room].Game.P1.Place(v.Coordinate, v.Piece, v.Direction)
					} else if message.Conn == roomList.M[message.Room].Participant {
						roomList.M[message.Room].Game.P2.Place(v.Coordinate, v.Piece, v.Direction)
					}
				}
			case util.ConnectionMsg:
				switch message.MsgType {
				case util.Create:
					roomList.CreateRoom(message.Conn)
				case util.Join:
					host, err := roomList.JoinRoom(message.Room, message.Conn)
					enc := gob.NewEncoder(message.Conn)
					if err != nil {
						util.ServerMsgEncoder(enc, util.StatusMsg{
							MsgType: util.Status,
							Action:  util.Join,
							Status:  false,
							Room:    -1,
						})
					}
					host_enc := gob.NewEncoder(host)
					util.ServerMsgEncoder(enc, util.StatusMsg{
						MsgType: util.Status,
						Action:  util.Join,
						Status:  true,
						Room:    message.Room,
					})
					util.ServerMsgEncoder(host_enc, util.StatusMsg{
						MsgType: util.Status,
						Action:  util.Join,
						Status:  true,
						Room:    message.Room,
					})
					// WARNING: Maybe move this to its own message
					roomList.M[message.Room].Game.PlayerTurn = game.PLAYER1
				case util.Delete:
					roomList.RemoveRoom(message.Room)
				case util.List:
					rooms := roomList.GetRooms()
					enc := gob.NewEncoder(message.Conn)
					util.ServerMsgEncoder(enc, util.RoomMsg{
						MsgType: util.List,
						Rooms:   rooms,
					})
				case util.Clear:
					roomList.ClearRooms()
				case util.Leave:
					host := roomList.M[message.Room].Host
					participant := roomList.M[message.Room].Participant
					hostEnc := gob.NewEncoder(host)
					partEnc := gob.NewEncoder(participant)
					util.ServerMsgEncoder(hostEnc, util.StatusMsg{
						MsgType: util.Exit,
					})
					util.ServerMsgEncoder(partEnc, util.StatusMsg{
						MsgType: util.Exit,
					})
					delete(roomList.M, message.Room)
				}
			}
		}
	}
}
