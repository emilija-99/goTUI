package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	appNameStyle    = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
	faint           = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)
	enumeratorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
)

func (m model) View() string {
	s := appNameStyle.Render("NOTES") + "\n\n"

	if m.state == titleView {
		s += "Note title: \n\n"
		s += m.textinput.View() + "\n\n"
		s += faint.Render("enter - save, esc - discard")
	}

	if m.state == bodyView {
		s += "Note body: \n\n"
		s += m.textarea.View() + "\n\n"
		s += faint.Render("ctrl+o - save, esc - discard")
	}

	if m.state == listView {
		for i, n := range m.notes {
			prefix := " "
			if i == m.listIndex {
				prefix = ">"
			}

			shortBody := strings.ReplaceAll(n.Body, "\n", "")

			if len(shortBody) > 30 {
				shortBody = shortBody[:30]
			}

			s += enumeratorStyle.Render(prefix) + n.Title + " | " + faint.Render(n.Body) + "\n\n"
		}
		s += faint.Render("n - new note, q - quit")
	}

	return s
}
