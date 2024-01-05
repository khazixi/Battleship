package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"

	// "log"
	// "net"
	// "bufio"
	// "os"
	// "strconv"
	// "strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khazixi/Battelship/ui"
	"github.com/khazixi/Battelship/util"
)

type ClientStore struct {
	connected bool
	activeID  int
}

func main() {

	gob.Register(util.ListAction{})
	gob.Register(util.JoinAction{})
	gob.Register(util.RoomMessage{})
	gob.Register(util.RoomsMessage{})
	gob.Register(util.CreateAction{})
	gob.Register(util.ConfirmationMessage{})

	// conn, err := net.Dial("tcp", ":8080")
	// if err != nil {
	// 	log.Fatalf("An error occured: %s\n", err)
	// 	return
	// }
	//
	// defer func() {
	// 	log.Println("Closing Connection")
	// 	conn.Close()
	// }()

	// var store ClientStore

	// fmt.Println(store)

	// reader := bufio.NewReader(os.Stdin)

	// go func() {
	// 	for {
	// 		decoder := gob.NewDecoder(conn)
	// 		message, err := util.MessageDecoder(decoder)
	//
	// 		if err != nil {
	// 			fmt.Println(err)
	// 			continue
	// 		}
	//
	// 		switch vm := message.(type) {
	// 		case util.RoomMessage:
	// 			if !store.connected {
	// 				store.activeID = vm.RoomID
	// 				store.connected = true
	//
	// 				fmt.Print(vm)
	// 			}
	// 		case util.ConfirmationMessage:
	// 			if !store.connected {
	// 				store.activeID = vm.RoomID
	// 				store.connected = true
	//
	// 				fmt.Print(vm)
	// 			}
	// 		}
	// 	}
	// }()

	// NOTE: Command Line Args should be used to pass the IP in the future
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

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

	// shell:
	//
	//	for {
	//		encoder := gob.NewEncoder(conn)
	//		fmt.Print(">>> ")
	//		readable, _ := reader.ReadString('\n')
	//
	//		fmt.Println(readable)
	//
	//		switch {
	//		case strings.HasPrefix(readable, "rooms"):
	//			util.ActionEncoder(encoder, util.ListAction{})
	//		case strings.HasPrefix(readable, "create"):
	//			if !store.connected {
	//				util.ActionEncoder(encoder, util.CreateAction{})
	//			}
	//		case strings.HasPrefix(readable, "join"):
	//			c := strings.Split(readable, " ")
	//			if len(c) != 2 {
	//				fmt.Println("Join takes 2 Arguments")
	//				continue
	//			}
	//			v, err := strconv.Atoi(strings.TrimSpace(c[1]))
	//			if err != nil {
	//				fmt.Println(err)
	//				continue
	//			}
	//			if !store.connected {
	//				util.ActionEncoder(encoder, util.JoinAction{RoomID: v})
	//			}
	//		case strings.HasPrefix(readable, "quit"):
	//			break shell
	//		default:
	//			fmt.Println("Not an action")
	//		}
	//	}
}
