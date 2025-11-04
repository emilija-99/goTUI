package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	listView  = 0
	titleView = 1
	bodyView  = 2
)

type model struct {
	state     uint
	store     *Store
	notes     []Note
	currNote  Note
	listIndex int
	textarea  textarea.Model
	textinput textinput.Model
}

func NewModel(store *Store) model {
	notes, err := store.GetNotes()
	if err != nil {
		log.Fatalf("unbale to get notes: %v", err)
	}

	return model{
		state:     listView,
		store:     store,
		notes:     notes,
		textarea:  textarea.New(),
		textinput: textinput.New(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// Only update active inputs depending on state
	switch m.state {
	case titleView:
		m.textinput, cmd = m.textinput.Update(msg)
		cmds = append(cmds, cmd)
	case bodyView:
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case listView:
			switch key {
			case "q":
				return m, tea.Quit

			case "n":
				m.textinput.SetValue("")
				m.textinput.Focus()
				m.textarea.Blur()
				m.currNote = Note{}
				m.state = titleView

			case "up", "k":
				if m.listIndex > 0 {
					m.listIndex--
				}

			case "down", "j":
				if m.listIndex < len(m.notes)-1 {
					m.listIndex++
				}

			case "enter":
				if len(m.notes) == 0 {
					break
				}
				m.currNote = m.notes[m.listIndex]
				m.textarea.SetValue(m.currNote.Body)
				m.textarea.Focus()
				m.textinput.Blur()
				m.state = bodyView
			}

		case titleView:
			switch key {
			case "enter":
				title := strings.TrimSpace(m.textinput.Value())
				if title != "" {
					m.currNote.Title = title
					m.textarea.SetValue("")
					m.textarea.Focus()
					m.textinput.Blur()
					m.state = bodyView
				}
			case "esc":
				m.textinput.Blur()
				m.state = listView
			}

		case bodyView:
			switch key {
			// Ctrl+S is reserved for “XOFF” — it pauses the terminal output (part of old “flow control”).
			case "ctrl+o":
				m.currNote.Body = m.textarea.Value()
				if err := m.store.SaveNote(m.currNote); err != nil {
					log.Printf("error saving note: %v", err)
				}

				var err error
				m.notes, err = m.store.GetNotes()
				if err != nil {
					log.Printf("error getting notes: %v", err)
				}

				m.currNote = Note{}
				m.textarea.Blur()
				m.state = listView

			case "esc":
				m.textarea.Blur()
				m.state = listView
			}
		}
	}

	return m, tea.Batch(cmds...)
}
