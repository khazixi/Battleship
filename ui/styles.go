package ui

import "github.com/charmbracelet/lipgloss"

var selected = lipgloss.
	NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("34"))

var border = lipgloss.
	NewStyle().
	PaddingLeft(2).
	PaddingRight(2).
	PaddingTop(1).
	PaddingBottom(1).
  Margin(1).
	BorderStyle(lipgloss.RoundedBorder())
