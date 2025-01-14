package bubbleforms

import (
	"fmt"
	"strings"
)

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
