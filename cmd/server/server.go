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
				Room:    roomID,
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
			log.Println("Recieved Action")
			msgch <- util.ConnectionMsg{
				MsgType: currentAction.Action,
				Room:    currentAction.Room,
				Conn:    conn,
			}

		case util.InitMsg:
			// TODO: Impliment game features
			msgch <- util.StartMsg{
				Conn:    conn,
				InitMsg: currentAction,
			}

		case util.PlaceMsg:
			msgch <- util.MarkMsg{
				Conn:     conn,
				PlaceMsg: currentAction,
			}
		}
	}
}

func gameLoop(msgch chan util.InternalMsg) {
	var roomList util.RoomList = util.MakeRoomList()
	for {
		select {
		case message := <-msgch:
			switch message := message.(type) {
			case util.StartMsg:
				for _, v := range message.Transmit {
					if message.Conn == roomList.M[message.Room].Host {
            err := roomList.M[message.Room].Game.P1.Place(v.Coordinate, v.Piece, v.Direction)
						if err != nil {
							enc := gob.NewEncoder(message.Conn)
							util.ServerMsgEncoder(enc, util.StatusMsg{
								MsgType: util.Status,
								Action:  util.Initialize,
								Status:  false,
								Room:    message.Room,
							})
						}
            roomList.M[message.Room].Game.P1.Reset()
					} else if message.Conn == roomList.M[message.Room].Participant {
            err := roomList.M[message.Room].Game.P2.Place(v.Coordinate, v.Piece, v.Direction)
						if err != nil {
							enc := gob.NewEncoder(message.Conn)
							util.ServerMsgEncoder(enc, util.StatusMsg{
								MsgType: util.Status,
								Action:  util.Initialize,
								Status:  false,
								Room:    message.Room,
							})
						}
            roomList.M[message.Room].Game.P2.Reset()
					}
				}

				enc := gob.NewEncoder(message.Conn)
				util.ServerMsgEncoder(enc, util.StatusMsg{
					MsgType: util.Status,
					Action:  util.Initialize,
					Status:  true,
					Room:    message.Room,
				})
			case util.MarkMsg:
				var hit bool
				if message.Conn == roomList.M[message.Room].Host {
					hit = roomList.M[message.Room].Game.P2.Mark(message.Mark)
				} else if message.Conn == roomList.M[message.Room].Participant {
					hit = roomList.M[message.Room].Game.P1.Mark(message.Mark)
				}
				switch {
				case roomList.M[message.Room].Game.P1.HasWin(): // TODO: Send the user that they won
					hostenc := gob.NewEncoder(roomList.M[message.Room].Host)
					partenc := gob.NewEncoder(roomList.M[message.Room].Participant)
					util.ServerMsgEncoder(hostenc, util.WinMsg{
						MsgType: util.GameResult,
						Winner:  game.PLAYER1,
					})
					util.ServerMsgEncoder(partenc, util.WinMsg{
						MsgType: util.GameResult,
						Winner:  game.PLAYER1,
					})
				case roomList.M[message.Room].Game.P2.HasWin(): // TODO: Send the user that they won
					hostenc := gob.NewEncoder(roomList.M[message.Room].Host)
					partenc := gob.NewEncoder(roomList.M[message.Room].Participant)
					util.ServerMsgEncoder(hostenc, util.WinMsg{
						MsgType: util.GameResult,
						Winner:  game.PLAYER2,
					})
					util.ServerMsgEncoder(partenc, util.WinMsg{
						MsgType: util.GameResult,
						Winner:  game.PLAYER2,
					})

				case hit: // TODO: Notify the users that a ship was hit
					hostenc := gob.NewEncoder(roomList.M[message.Room].Host)
					partenc := gob.NewEncoder(roomList.M[message.Room].Participant)
					util.ServerMsgEncoder(hostenc, util.HitMsg{
						MsgType:    util.Hit,
						Hit:        true,
						Coordinate: message.Mark,
					})
					util.ServerMsgEncoder(partenc, util.HitMsg{
						MsgType:    util.Hit,
						Hit:        true,
						Coordinate: message.Mark,
					})
				default: // TODO: Notify the players that a ship was not hit
					hostenc := gob.NewEncoder(roomList.M[message.Room].Host)
					partenc := gob.NewEncoder(roomList.M[message.Room].Participant)
					util.ServerMsgEncoder(hostenc, util.HitMsg{
						MsgType:    util.Hit,
						Hit:        false,
						Coordinate: message.Mark,
					})
					util.ServerMsgEncoder(partenc, util.HitMsg{
						MsgType:    util.Hit,
						Hit:        false,
						Coordinate: message.Mark,
					})
				}
				roomList.M[message.Room].Game.P2.HasWin()

			case util.ConnectionMsg:
				switch message.MsgType {
				case util.Create:
					room := roomList.CreateRoom(message.Conn)
					enc := gob.NewEncoder(message.Conn)
					util.ServerMsgEncoder(enc, util.StatusMsg{
						MsgType: util.Status,
						Action:  util.Create,
						Status:  true,
						Room:    room,
					})

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
				case util.Delete:
					roomList.RemoveRoom(message.Room)
					enc := gob.NewEncoder(message.Conn)
					util.ServerMsgEncoder(enc, util.StatusMsg{
						MsgType: util.Status,
						Action:  util.Delete,
						Status:  true,
						Room:    message.Room,
					})
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
