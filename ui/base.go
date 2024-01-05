package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/khazixi/Battelship/util"
)

type Base struct {
	cursor  int
	choices []string
	Options
}

func MakeModel(option Options) Base {
	return Base{
		cursor:  0,
		choices: []string{"New Game", "Join Game", "Exit"},
		Options: option,
	}
}

func (m Base) Init() tea.Cmd {
	return tea.Batch(
		m.Options.Listen(m.Options.msgch),
		m.Options.Process(m.Options.msgch),
    m.Options.Send(util.CreateAction{}),
	)
}

func (m Base) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			switch m.choices[m.cursor] {
			case "Exit":
				return m, tea.Quit
			case "New Game":
				return RoomCreatorModel{Options: m.Options}, m.Options.Send(util.CreateAction{})
			case "Join Game":
				return RoomViewModel{Options: m.Options}, m.Options.Send(util.ListAction{})
			}
		}
	}

	return m, nil
}

func (m Base) View() string {
	s := "The Score\n\n"
	for i, c := range m.choices {
		if i == m.cursor {
			s += fmt.Sprintf("> %s\n", selected.Render(c))
		} else {
			s += fmt.Sprintf("  %s\n", c)
		}
	}
	return border.Render(s)
}
