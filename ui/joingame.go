package ui

import tea "github.com/charmbracelet/bubbletea"

type RoomViewModel struct{
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
	return border.Render("Rooms")
}
