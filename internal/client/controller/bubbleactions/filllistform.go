package bubbleactions

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	ifcs "github.com/fasdalf/train-go-gophkeeper/internal/client/interfaces"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
	"log/slog"
)

var typeCheckFillListFormService ListFormFiller = &FillListFormService{}
var typeCheckFillListFormMock ListFormFiller = &FillListFormMock{}

type ListFormFiller interface {
	Fill(form *bubbleforms.InputsModel) tea.Model
}

type FillListFormService struct {
	syncroniser ifcs.Syncroniser
	repository  ifcs.LocalRepository
}

func NewFillListFormService(
	syncroniser ifcs.Syncroniser,
	repository ifcs.LocalRepository,
) *FillListFormService {
	return &FillListFormService{
		syncroniser: syncroniser,
		repository:  repository,
	}
}

func (s *FillListFormService) Fill(form *bubbleforms.InputsModel) tea.Model {
	m := tea.Model(form)
	listModel, ok := (*form.ParentModel).(*bubbleforms.ListModel)
	if !ok || listModel == nil {
		slog.Error("could not cast model to listModel")
		return &bubbleforms.MessageModel{Message: "could not cast model to listModel", Next: &m}
	}

	if err := s.syncroniser.Sync(); err != nil {
		slog.Error("could not sync data", "err", err)
		return &bubbleforms.MessageModel{Message: err.Error(), Next: &m}
	}

	items := []list.Item{}
	for i, name := range s.repository.ListSecrets() {
		items = append(items, bubbleforms.ListItem{
			Name: name,
			Desc: "",
			ID:   i,
		})
	}
	listModel.List.SetItems(items)

	return listModel
}

type FillListFormMock struct {
	Result tea.Model
}

func (s *FillListFormMock) Fill(form *bubbleforms.InputsModel) tea.Model {
	return s.Result
}
