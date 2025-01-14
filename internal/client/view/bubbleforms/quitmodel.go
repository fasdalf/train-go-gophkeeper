package bubbleforms

import (
	tea "github.com/charmbracelet/bubbletea"
	"log/slog"
	"reflect"
)

type QuitModel struct {
}

func (m *QuitModel) Init() tea.Cmd {
	return tea.Quit
}

func (m *QuitModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Info("QuitModel update called", "msg", msg, "m", m, "type", reflect.TypeOf(m), "type", reflect.TypeOf(msg))
	return m, tea.Quit
}

func (m *QuitModel) View() string {
	return ""
}

type QuitController struct {
}

func (c *QuitController) Handle(tea.Model) tea.Model {
	return &QuitModel{}
}
