package bubbleforms

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	ifcs "github.com/fasdalf/train-go-gophkeeper/internal/client/interfaces"
	"log/slog"
)

const listTitle = "secrets list"

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ListItem struct {
	Name, Desc string
	ID         int
}

func (i ListItem) Title() string       { return i.Name }
func (i ListItem) Description() string { return i.Desc }
func (i ListItem) FilterValue() string { return i.Name }

type ListModel struct {
	List          list.Model
	AddHandler    ifcs.ButtonHandler
	EditHandler   ifcs.ButtonHandler
	DeleteHandler ifcs.ButtonHandler
	ReloadHandler ifcs.ButtonHandler
}

func (m *ListModel) Init() tea.Cmd {
	return forceWindowSize
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Info("ListModel update called", "msg", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			slog.Info("ListModel got ctrl+c")
			return m, tea.Quit
		case "insert":
			m2 := m.AddHandler.Handle(m)
			return m2, m2.Init()
		case "enter":
			m2 := m.EditHandler.Handle(m)
			return m2, m2.Init()
		case "delete":
			//m2 := m.DeleteHandler.Handle(m)
			//return m2, m2.Init()
			return m, nil
		case "ctrl+r":
			//m2 := m.ReloadHandler.Handle(m)
			//return m2, m2.Init()
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m *ListModel) View() string {
	return docStyle.Render(m.List.View())
}

func NewListModel(addHandler ifcs.ButtonHandler, editHandler ifcs.ButtonHandler) *ListModel {
	lid := list.NewDefaultDelegate()
	lid.ShowDescription = false
	l := list.New([]list.Item{}, lid, 0, 0)
	l.Title = listTitle
	l.SetShowTitle(true)
	// TODO: ##@@ implement AdditionalShortHelpKeys + AdditionalFullHelpKeys
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	m := ListModel{
		List:        l,
		AddHandler:  addHandler,
		EditHandler: editHandler,
	}

	return &m
}
