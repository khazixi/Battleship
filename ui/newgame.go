package ui

import tea "github.com/charmbracelet/bubbletea"

type RoomCreatorModel struct {
  Options
}

func MakeRoomCreator() RoomCreatorModel {
 return RoomCreatorModel{}
}

func (m RoomCreatorModel) Init() tea.Cmd {
  return nil
}


func (m RoomCreatorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (r RoomCreatorModel) View() string {
  return border.Render("hi")
}
