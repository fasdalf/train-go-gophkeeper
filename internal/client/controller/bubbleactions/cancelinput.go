package bubbleactions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
	"log/slog"
)

var typeCheckCancelInputController bubbleforms.ButtonController = &CancelInputController{}

type CancelInputController struct {
}

func NewCancelInputController() *CancelInputController {
	return &CancelInputController{}
}

func (s *CancelInputController) Handle(m tea.Model) tea.Model {
	form, ok := m.(*bubbleforms.InputsModel)
	if !ok || form == nil || form.ParentModel == nil {
		slog.Error("nowhere to return")
		return &bubbleforms.MessageModel{Message: "nowhere to return", Next: &m}
	}

	return *form.ParentModel
}
