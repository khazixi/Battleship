package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khazixi/Battelship/ui"
	"github.com/khazixi/Battelship/util"
)

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
}
