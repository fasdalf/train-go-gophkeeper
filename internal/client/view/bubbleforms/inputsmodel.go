package bubbleforms

import (
	"fmt"
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	ifcs "github.com/fasdalf/train-go-gophkeeper/internal/client/interfaces"
	"log/slog"
	"strings"
)

const (
	InputTextPlain = iota
	InputTextPassword
	InputTextarea
)

const (
	promptActive   = " │ "
	promptInactive = " - "
)

type InputType = int
type InputField struct {
	Type  InputType
	Title string
	Value string
}
type InputButton struct {
	Label   string
	Handler ifcs.ButtonHandler
}

type InputsModel struct {
	ParentModel *tea.Model
	Fields      []InputField
	focusIndex  int
	buttons     []InputButton
	textinputs  []*textinput.Model
}

func NewInputsModel(fields []InputField, buttons []InputButton, ParentModel *tea.Model) *InputsModel {
	m := InputsModel{
		ParentModel: ParentModel,
		focusIndex:  0,
		Fields:      fields,
		buttons:     buttons,
		textinputs:  make([]*textinput.Model, len(fields)),
	}

	for i, f := range m.Fields {
		switch f.Type {
		case InputTextPlain:
		case InputTextPassword:
		default:
			continue
		}

		t := textinput.New()
		t.Cursor.SetMode(cursor.CursorStatic)

		t.Prompt = promptInactive
		if i == 0 {
			t.Prompt = promptActive
			t.Focus()
		}

		if f.Type == InputTextPassword {
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		t.SetValue(f.Value)
		m.textinputs[i] = &t
	}

	return &m
}

func (m *InputsModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *InputsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Info("InputsModel update called", "msg", msg, "m", m)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			slog.Info("InputsModel got ctrl+c")
			return m, tea.Quit
		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the button was focused?
			if s == "enter" && m.focusIndex >= len(m.Fields) {
				buttonHandler := m.buttons[m.focusIndex-len(m.Fields)].Handler
				newModel := m.callHandler(buttonHandler)
				return newModel, newModel.Init()
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
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
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
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

func (m *InputsModel) callHandler(buttonHandler ifcs.ButtonHandler) tea.Model {
	for i := range m.Fields {
		switch m.Fields[i].Type {
		case InputTextPlain, InputTextPassword:
			m.Fields[i].Value = m.textinputs[i].Value()
		}
	}

	return buttonHandler.Handle(m)
}

func (m *InputsModel) View() string {
	var b strings.Builder

	for i, f := range m.Fields {
		switch f.Type {
		case InputTextPlain, InputTextPassword:
			b.WriteString(f.Title)
			b.WriteRune('\n')
			b.WriteString(m.textinputs[i].View())
			b.WriteString("\n\n")
		}
	}

	fl := len(m.Fields)
	for i, btn := range m.buttons {
		fmt.Fprintf(&b, "%s\n", renderButton(btn.Label, m.focusIndex == (fl+i)))
		if i < len(m.buttons)-1 {
			b.WriteRune('\n')
		}
	}

	return b.String()
}

func renderButton(label string, focused bool) string {
	prompt := promptInactive
	if focused {
		prompt = promptActive
	}
	return fmt.Sprintf("%s [ %s ]", prompt, label)
}
