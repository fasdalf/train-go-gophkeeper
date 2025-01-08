package bubbleforms

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/term"
	"github.com/muesli/reflow/wordwrap"
	"log/slog"
	"os"
	"reflect"
)

type MessageModel struct {
	Message string
	Next    *tea.Model
	Width   int
	Height  int
}

func forceWindowSize() tea.Msg {
	m := tea.WindowSizeMsg{}
	m.Width, m.Height, _ = term.GetSize(os.Stdout.Fd())
	return m
}

func (m *MessageModel) Init() tea.Cmd {
	return forceWindowSize
}

func (m *MessageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Info("MessageModel update called", "msg", msg, "m", m, "type", reflect.TypeOf(m), "type", reflect.TypeOf(msg))
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height // height is available too
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			slog.Info("MessageModel got ctrl+c")
			return m, tea.Quit
		default:
			if m.Next != nil {
				return *m.Next, (*m.Next).Init()
			}
			return m, tea.Quit
		}
	default:
		return m, nil
	}
}

func (m *MessageModel) View() string {
	action := "quit"
	if m.Next != nil {
		action = "continue"
	}
	return wordwrap.String(fmt.Sprintf("%s\n\nPress ANY key to %s.", m.Message, action), m.Width)
}
