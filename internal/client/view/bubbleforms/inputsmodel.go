package bubbleforms

import (
	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	InputTextPlain = iota
	InputTextPassword
	InputTextarea
)

const (
	promptActive          = " │ "
	promptInactive        = " - "
	inputsModelKBNextKey1 = "tab"
	inputsModelKBBackKey1 = "shift+tab"
	inputsModelKBNextKey2 = "down"
	inputsModelKBBackKey2 = "up"
	inputsModelKBEnterKey = "enter"
)

type InputType = int
type InputField struct {
	Type  InputType
	Title string
	Value string
}
type InputButton struct {
	Label   string
	Handler ButtonController
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
