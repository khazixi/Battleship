package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khazixi/Battelship/game"
	"github.com/khazixi/Battelship/ui"
	"github.com/khazixi/Battelship/util"
)

type Store struct {
	room   int
	conn   net.Conn
	placed []game.Transmit
	player game.Turn
	board  game.Board
}

func main() {
	gob.Register(util.StatusMsg{})
	gob.Register(util.RoomMsg{})
	gob.Register(util.ActionMsg{})
	gob.Register(util.InitMsg{})

	// NOTE: Command Line Args should be used to pass the IP in the future

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	isTui := flag.Bool("tui", false, "Using the tui option toggles the Tui")
	flag.Parse()

	if *isTui {
		program := tea.NewProgram(ui.MakeModel(ui.MakeOption(conn)))

		if len(os.Getenv("DEBUG")) > 0 {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				fmt.Println("fatal:", err)
				os.Exit(1)
			}
			defer f.Close()
		}

		_, err = program.Run()
		if err != nil {
			fmt.Println("Failed to Create the UI")
			os.Exit(1)
		}
	} else {
		store := Store{conn: conn}
		fmt.Print(">>> ")
		scnr := bufio.NewScanner(os.Stdin)
		for scnr.Scan() {
			str := scnr.Text()
			command := strings.Fields(str)
			switch command[0] {
			case "create":
				err = store.createCmd(command)
				if err != nil {
					fmt.Println(err)
				} else {
					store.processServerMsg()
				}
			case "join":
				err = store.joinCmd(command)
				if err != nil {
					fmt.Println(err)
				} else {
					store.processServerMsg()
				}
			case "mark":
				err = store.markCmd(command)
				if err != nil {
					fmt.Println(err)
				} else {
					store.processServerMsg()
				}
			case "place":
				err = store.placeCmd(command)
				if err != nil {
					fmt.Println(err)
				} else {
					store.processServerMsg()
				}
			case "list":
				err = store.listCmd(command)
				if err != nil {
					fmt.Println(err)
				} else {
					store.processServerMsg()
				}
			case "help":
				fmt.Println("Usage: create")
				fmt.Println("Usage: join <room>")
				fmt.Println("Usage: mark <coordinate>")
				fmt.Println("Usage: place <piece> <direction> <coordinate>")
			default:
				fmt.Println(command[0], " is not a command")
			}
			fmt.Print(">>> ")
		}
		fmt.Println("This should not be a TUI")
	}

}

func (s Store) listCmd(command []string) error {
	if len(command) != 1 {
		return errors.New("Usage: list")
	}

	encoder := gob.NewEncoder(s.conn)
	util.ClientMsgEncoder(encoder, util.ActionMsg{
		MsgType: util.Action,
		Action:  util.List,
		Room:    -1,
	})

	return nil
}

func (s Store) createCmd(command []string) error {
	if len(command) != 1 {
		return errors.New("Usage: create")
	}
	encoder := gob.NewEncoder(s.conn)
	util.ClientMsgEncoder(encoder, util.ActionMsg{
		MsgType: util.Action,
		Action:  util.Create,
		Room:    -1,
	})

	return nil
}

func (s Store) joinCmd(command []string) error {
	if len(command) != 2 {
		return errors.New("Usage: join <room>")
	}
	room, err := strconv.ParseInt(command[1], 10, 64)
	if err != nil {
		fmt.Println(err)
		return err
	}
	encoder := gob.NewEncoder(s.conn)
	util.ClientMsgEncoder(encoder, util.ActionMsg{
		MsgType: util.Action,
		Action:  util.Join,
		Room:    int(room),
	})

	return nil
}

func (s Store) markCmd(command []string) error {
	if len(command) != 2 {
		return errors.New("Usage: mark <coordinate>")
	}

	instruction, err := game.ParseInstruction(command[1])
	if err != nil {
		return errors.New("Failed to properly parse the instruction")
	}

	encoder := gob.NewEncoder(s.conn)
	util.ClientMsgEncoder(encoder, util.PlaceMsg{
		MsgType: util.Mark,
		Room:    s.room,
		Mark:    instruction,
	})

	return nil
}

func (s *Store) placeCmd(command []string) error {
	var piece game.Piece
	var direction game.Direction
	if len(command) != 4 {
		return errors.New("Usage: place <piece> <direction> <coordinate>")
	}

	switch command[1] {
	case "carrier":
		piece = game.CARRIER
	case "destroyer":
		piece = game.DESTROYER
	case "patrolboat":
		piece = game.PATROLBOAT
	case "submarine":
		piece = game.SUBMARINE
	case "battleship":
		piece = game.BATTLESHIP
	default:
		return errors.New("Incorrect Ship\nMust be either: carrier, destroyer, patrolboat, submarine, battleship")
	}

	switch command[2] {
	case "left":
		direction = game.LEFT
	case "up":
		direction = game.UP
	case "right":
		direction = game.RIGHT
	case "down":
		direction = game.DOWN
	default:
		return errors.New("Incorrect Direction\nMust be either: left, up, right, down")
	}

	instruction, err := game.ParseInstruction(command[3])
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = s.board.Place(instruction, piece, direction)
	if err != nil {
		fmt.Println(err)
		return err
	}

	s.placed = append(s.placed, game.Transmit{Coordinate: instruction, Piece: piece, Direction: direction})

	if len(s.placed) == 5 && s.room != -1 {
		encoder := gob.NewEncoder(s.conn)
		util.ClientMsgEncoder(encoder, util.InitMsg{
			MsgType:  util.Initialize,
			Room:     s.room,
			Transmit: [5]game.Transmit(s.placed),
		})
	}

	return nil
}

func (s *Store) processServerMsg() {
	decoder := gob.NewDecoder(s.conn)
	msg, err := util.ServerMsgDecoder(decoder)
  fmt.Println(msg)
	if err != nil {
		fmt.Println(err)
	}
	switch msg := msg.(type) {
	case util.StatusMsg:
		if msg.Status && msg.Room != -1 {
			s.room = msg.Room
		}
	case util.RoomMsg:
		for _, room := range msg.Rooms {
			fmt.Println("Room |", room)
		}
	case util.WinMsg:
		fmt.Println("won")
	case util.HitMsg:
		if msg.Hit {
			fmt.Println("You were hit")
		} else {
			fmt.Println("You were not hit")
		}
	case util.InitializedMsg:
		s.player = msg.Player
	}
}
