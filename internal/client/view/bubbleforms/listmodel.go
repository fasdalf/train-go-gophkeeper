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
	slog.Info("##@@ update called", "msg", msg)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			slog.Info("##@@ got ctrl+c")
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
	// TODO: ##@@ cleanup
	//items := []list.Item{
	//	//ListItem{Name: "Raspberry Pi’s", Desc: "I have ’em all over my house"},
	//	//ListItem{Name: "Nutella", Desc: "It's good on toast"},
	//	//ListItem{Name: "Bitter melon", Desc: "It cools you down"},
	//	//ListItem{Name: "Nice socks", Desc: "And by that I mean socks without holes"},
	//	//ListItem{Name: "Eight hours of sleep", Desc: "I had this once"},
	//	//ListItem{Name: "Cats", Desc: "Usually"},
	//	//ListItem{Name: "Plantasia, the album", Desc: "My plants love it too"},
	//	//ListItem{Name: "Pour over coffee", Desc: "It takes forever to make though"},
	//	//ListItem{Name: "VR", Desc: "Virtual reality...what is there to say?"},
	//	//ListItem{Name: "Noguchi Lamps", Desc: "Such pleasing organic forms"},
	//	//ListItem{Name: "Linux", Desc: "Pretty much the best OS"},
	//	//ListItem{Name: "Business school", Desc: "Just kidding"},
	//	//ListItem{Name: "Pottery", Desc: "Wet clay is a great feeling"},
	//	//ListItem{Name: "Shampoo", Desc: "Nothing like clean hair"},
	//	//ListItem{Name: "Table tennis", Desc: "It’s surprisingly exhausting"},
	//	//ListItem{Name: "Milk crates", Desc: "Great for packing in your extra stuff"},
	//	//ListItem{Name: "Afternoon tea", Desc: "Especially the tea sandwich part"},
	//	//ListItem{Name: "Stickers", Desc: "The thicker the vinyl the better"},
	//	//ListItem{Name: "20° Weather", Desc: "Celsius, not Fahrenheit"},
	//	//ListItem{Name: "Warm light", Desc: "Like around 2700 Kelvin"},
	//	//ListItem{Name: "The vernal equinox", Desc: "The autumnal equinox is pretty good too"},
	//	//ListItem{Name: "Gaffer’s tape", Desc: "Basically sticky fabric"},
	//	//ListItem{Name: "Terrycloth", Desc: "In other words, towel fabric"},
	//}

	lid := list.NewDefaultDelegate()
	lid.ShowDescription = false
	l := list.New([]list.Item{}, lid, 0, 0)
	l.Title = listTitle
	l.SetShowTitle(true)
	// TODO: implement AdditionalShortHelpKeys + AdditionalFullHelpKeys
	l.SetShowHelp(false)
	l.DisableQuitKeybindings()
	m := ListModel{
		List:        l,
		AddHandler:  addHandler,
		EditHandler: editHandler,
	}

	return &m
}
