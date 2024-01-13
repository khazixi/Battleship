package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khazixi/Battelship/util"
)

type RoomCreatorModel struct {
	Options
	rooms []int
}

func (m RoomCreatorModel) Init() tea.Cmd {
	return tea.Batch(
		m.Options.Listen(m.Options.msgch),
		m.Options.Process(m.Options.msgch),
	)
}

func (m RoomCreatorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case util.StatusMsg:
		fmt.Println(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
      m.Options.Room = -1
			return MakeModel(m.Options), m.Send(util.ActionMsg{MsgType: util.Action, Action: util.Leave})
		}
	}

	return m, nil
}

func (r RoomCreatorModel) View() string {

	return border.Render("hi")
}
