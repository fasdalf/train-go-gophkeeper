package bubbleforms

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log/slog"
)

const listTitle = "secrets list"
const listModelKBAddKey = "insert"
const listModelKBAddDesc = "add secret"
const listModelKBEditKey = "enter"
const listModelKBEditDesc = "edit current secret"

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type ListItem struct {
	Name, Desc string
	ID         int
}

func (i ListItem) Title() string       { return i.Name }
func (i ListItem) Description() string { return i.Desc }
func (i ListItem) FilterValue() string { return i.Name }

type KeyController struct {
	Controller ButtonController
	KeyBinding key.Binding
}

type ListModel struct {
	List           list.Model
	KeyControllers []KeyController
}

func (m *ListModel) Init() tea.Cmd {
	return forceWindowSize
}

func (m *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	slog.Info("ListModel update called", "msg", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		for _, kc := range m.KeyControllers {
			if key.Matches(msg, kc.KeyBinding) {
				m2 := kc.Controller.Handle(m)
				return m2, m2.Init()
			}
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

func NewListModel(addHandler ButtonController, editHandler ButtonController) *ListModel {
	quitController := KeyController{
		Controller: &QuitController{},
		KeyBinding: key.NewBinding(
			key.WithKeys(KBQuitKey),
			key.WithHelp(KBQuitKey, KBQuitDesc),
		),
	}
	addController := KeyController{
		Controller: addHandler,
		KeyBinding: key.NewBinding(
			key.WithKeys(listModelKBAddKey),
			key.WithHelp(listModelKBAddKey, listModelKBAddDesc),
		),
	}
	editController := KeyController{
		Controller: editHandler,
		KeyBinding: key.NewBinding(
			key.WithKeys(listModelKBEditKey),
			key.WithHelp(listModelKBEditKey, listModelKBEditDesc),
		),
	}

	lid := list.NewDefaultDelegate()
	lid.ShowDescription = false
	l := list.New([]list.Item{}, lid, 0, 0)
	l.Title = listTitle
	l.SetShowTitle(true)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{quitController.KeyBinding, addController.KeyBinding, editController.KeyBinding}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys
	l.SetShowHelp(true)
	l.DisableQuitKeybindings()
	m := ListModel{
		List:           l,
		KeyControllers: []KeyController{quitController, addController, editController},
	}

	return &m
}
