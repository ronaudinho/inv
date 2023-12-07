package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	survival "github.com/ronaudinho/inv/cmd/survival"
)

type model struct {
	Choice     int
	Chosen     bool
	Loaded     bool
	Quitting   bool
	modeChosen int
}

// Init implements tea.Model.
func (model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.Quitting = true
			return m, tea.Quit
		}
	}

	if !m.Chosen {
		return updateChoice(msg, m)
	}

	return m, nil
}

// View implements tea.Model.
func (m model) View() string {
	var s string
	if m.Quitting {
		return "\n  See you later!\n\n"
	}
	if !m.Chosen {
		s = choicesView(m)
	} else {
		s = chosenView(m)
	}
	return s
}

func main() {
	initialModel := model{0, false, false, false, 0}
	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		fmt.Println("could not start program:", err)
	}
}

func selectChoise(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.Choice++
			if m.Choice > 5 {
				m.Choice = 5
			}
		case "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			m.modeChosen = m.Choice
			switch m.modeChosen {
			case 0:
				// survival
				survival.Survival()
			case 1:
				// classic
			case 2:
				// combo
			case 3:
				// hattrick
			case 4:
				// random
			case 5:
				// endless

			}
		}
	}

	return m, nil
}

func checkbox(label string, checked bool) string {
	if checked {
		return "[x] " + label
	}
	return fmt.Sprintf("[ ] %s", label)
}

func choicesView(m model) string {
	c := m.Choice
	tpl := "Select Game Mode?\n\n"
	tpl += "%s\n\n"
	tpl += "up/down: select <-> enter: choose <-> q: quit"

	choices := fmt.Sprintf(
		"%s\n%s\n%s\n%s\n%s\n%s",
		checkbox("Survival", c == 0),
		checkbox("Classic", c == 1),
		checkbox("Combo", c == 2),
		checkbox("Hattrick", c == 3),
		checkbox("Random", c == 4),
		checkbox("Endless", c == 5),
	)

	return fmt.Sprintf(tpl, choices)
}

func chosenView(m model) string {
	var msg string

	return msg
}

func updateChoice(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "down":
			m.Choice++
			if m.Choice > 5 {
				m.Choice = 5
			}
		case "up":
			m.Choice--
			if m.Choice < 0 {
				m.Choice = 0
			}
		case "enter":
			m.Chosen = true
			m.modeChosen = m.Choice
			// this is work but not like what I want. I dont know how i can running survival mode after I choose it from game mode option
			switch m.modeChosen {
			case 0:
				// survival
				survival.Survival()
			case 1:
				// classic
			case 2:
				// combo
			case 3:
				// hattrick
			case 4:
				// random
			case 5:
				// endless

			}
			return m, tea.Quit
		}
	}

	return m, nil
}
