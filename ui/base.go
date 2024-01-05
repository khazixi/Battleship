package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
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
				return MakeRoomCreator(), nil
			case "Join Game":
        return RoomViewModel{Options: m.Options}, nil
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
	// s += fmt.Sprintf("Cursor Position: [%d]\n\n", m.cursor)
	return border.Render(s)
}
