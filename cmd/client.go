package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

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
	gob.Register(util.CreateAction{})
  gob.Register(util.ConfirmationMessage{})

	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalf("An error occured: %s\n", err)
		return
	}

	defer func() {
		log.Println("Closing Connection")
		conn.Close()
	}()

	var store ClientStore

  fmt.Println(store)

	reader := bufio.NewReader(os.Stdin)

	go func() {
		for {
			decoder := gob.NewDecoder(conn)
			message, err := messageDecoder(decoder)

			if err != nil {
        fmt.Println(err)
				continue
			}

			switch vm := message.(type) {
			case util.RoomMessage:
				if !store.connected {
					store.activeID = vm.RoomID
					store.connected = true

          fmt.Print(vm)
				}
      case util.ConfirmationMessage:
				if !store.connected {
					store.activeID = vm.RoomID
					store.connected = true

          fmt.Print(vm)
				}
			}
		}
	}()


shell:
	for {
		encoder := gob.NewEncoder(conn)
		fmt.Print(">>> ")
		readable, _ := reader.ReadString('\n')

		fmt.Println(readable)

		switch {
		case strings.HasPrefix(readable, "rooms"):
			actionEncoder(encoder, util.CreateListAction(1))
		case strings.HasPrefix(readable, "create"):
			if !store.connected {
				actionEncoder(encoder, util.CreateCreateAction(1))
			}
		case strings.HasPrefix(readable, "join"):
			c := strings.Split(readable, " ")
			if len(c) != 2 {
				fmt.Println("Join takes 2 Arguments")
				continue
			}
			v, err := strconv.Atoi(strings.TrimSpace(c[1]))
			if err != nil {
				fmt.Println(err)
				continue
			}
			if !store.connected {
				actionEncoder(encoder, util.CreateJoinAction(1, v))
			}
		case strings.HasPrefix(readable, "quit"):
			break shell
		default:
			fmt.Println("Not an action")
		}
	}
}

func actionEncoder(enc *gob.Encoder, a util.Action) {
	err := enc.Encode(&a)
	if err != nil {
		log.Fatal("Failed to encode", err)
	}
}

func messageDecoder(dec *gob.Decoder) (util.Message, error) {
	var message util.Message
	err := dec.Decode(&message)
	return message, err
}
