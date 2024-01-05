package ui

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khazixi/Battelship/util"
)

type RoomViewModel struct{
  rooms []int
  Options
}

func (m RoomViewModel) Init() tea.Cmd {
	return tea.Batch(
    m.Options.Listen(m.Options.msgch),
    m.Options.Process(m.Options.msgch),
  )
}

func (m RoomViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
  case util.Message:
    log.Println(msg)
		switch msg := msg.(type) {
		case util.RoomsMessage:
      m.rooms = msg.Rooms
      return m, nil
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			return MakeModel(m.Options), nil
		}
	}

	return m, nil
}

func (r RoomViewModel) View() string {
  if len(r.rooms) == 0 {
    return "No Rooms Available to join"
  }

  s := ""
  for _, room := range r.rooms {
    s += fmt.Sprintf("Room | %d\n", room)
  }

	return border.Render(s)
}
