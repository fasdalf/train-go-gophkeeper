package bubbleforms

import (
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
)

func (m *InputsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Info("InputsModel update called", "msg", msg, "m", m)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch s := msg.String(); s {
		case KBQuitKey:
			slog.Info("InputsModel got ctrl+c")
			return m, tea.Quit
		// Call button controller
		case inputsModelKBEnterKey:
			// Did the user press enter while the button was focused?
			if s == inputsModelKBEnterKey && m.focusIndex >= len(m.Fields) {
				buttonHandler := m.buttons[m.focusIndex-len(m.Fields)].Handler
				newModel := m.callHandler(buttonHandler)
				return newModel, newModel.Init()
			}
			fallthrough
		// Set focus to next input
		case inputsModelKBNextKey1, inputsModelKBNextKey2, inputsModelKBBackKey1, inputsModelKBBackKey2:
			return m.navigateInputs(s)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *InputsModel) navigateInputs(s string) (tea.Model, tea.Cmd) {
	// Cycle indexes
	if s == inputsModelKBBackKey1 || s == inputsModelKBBackKey2 {
		m.focusIndex--
	} else {
		m.focusIndex++
	}

	lenAll := len(m.Fields) + len(m.buttons)
	if m.focusIndex >= lenAll {
		m.focusIndex = 0
	} else if m.focusIndex < 0 {
		m.focusIndex = lenAll - 1
	}

	cmds := make([]tea.Cmd, len(m.textinputs))
	for i := 0; i <= len(m.textinputs)-1; i++ {
		if m.textinputs[i] == nil {
			continue
		}

		if i == m.focusIndex {
			m.textinputs[i].Prompt = promptActive
			// Set focused state
			cmds[i] = m.textinputs[i].Focus()
			continue
		}
		m.textinputs[i].Prompt = promptInactive
		// Remove focused state
		m.textinputs[i].Blur()
	}

	return m, tea.Batch(cmds...)
}

func (m *InputsModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.textinputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.textinputs {
		if m.textinputs[i] == nil {
			continue
		}
		ti, cmd := m.textinputs[i].Update(msg)
		m.textinputs[i], cmds[i] = &ti, cmd
	}

	return tea.Batch(cmds...)
}

func (m *InputsModel) callHandler(buttonHandler ButtonController) tea.Model {
	for i := range m.Fields {
		switch m.Fields[i].Type {
		case InputTextPlain, InputTextPassword:
			m.Fields[i].Value = m.textinputs[i].Value()
		}
	}

	return buttonHandler.Handle(m)
}
