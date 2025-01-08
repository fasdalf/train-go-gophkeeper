package bubbleactions

import (
	tea "github.com/charmbracelet/bubbletea"
	ifcs "github.com/fasdalf/train-go-gophkeeper/internal/client/interfaces"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/convertors"
	"log/slog"
)

var typeCheckAddSecretController ifcs.ButtonHandler = &AddSecretController{}

type AddSecretController struct {
	service         ifcs.RemoteService
	syncroniser     ifcs.Syncroniser
	repository      ifcs.LocalRepository
	listFormFiller  ListFormFiller
	secretDataValue any
}

func NewAddSecretController(
	service ifcs.RemoteService,
	syncroniser ifcs.Syncroniser,
	repository ifcs.LocalRepository,
	listFormFiller ListFormFiller,
) *AddSecretController {
	return &AddSecretController{
		service:        service,
		syncroniser:    syncroniser,
		repository:     repository,
		listFormFiller: listFormFiller,
	}
}

func (s *AddSecretController) setSecretDataValue(v any) {
	s.secretDataValue = v
}

func (s *AddSecretController) Handle(m tea.Model) tea.Model {
	form, ok := m.(*bubbleforms.InputsModel)
	if !ok || form == nil || len(form.Fields) == 0 {
		slog.Error("not an InputsModel or empty Fields")
		return &bubbleforms.MessageModel{Message: "not an InputsModel or empty Fields", Next: &m}
	}
	listForm, ok := (*form.ParentModel).(*bubbleforms.ListModel)
	if !ok || listForm == nil {
		slog.Error("not a ListModel")
		return &bubbleforms.MessageModel{Message: "not a ListModel", Next: &m}
	}
	listFormTea := tea.Model(listForm)
	secret := entity.Secret{}
	sData := entity.SecretData{}
	switch s.secretDataValue.(type) {
	case entity.Password:
		sData = convertors.FieldsToPassword(form.Fields, nil)
	default:
		return &bubbleforms.MessageModel{Message: "invalid secret type", Next: &listFormTea}
	}
	err := secret.SetData(sData)
	if err != nil {
		slog.Error("can't encode secret", "error", err)
		return &bubbleforms.MessageModel{Message: "can't encode secret", Next: &listFormTea}
	}
	secret.ServerID, secret.ServerUpdatedAt, err = s.syncroniser.Add(secret)
	if err != nil {
		slog.Error("can't send secret", "error", err)
		return &bubbleforms.MessageModel{Message: "can't send secret", Next: &m}
	}
	return s.listFormFiller.Fill(form)
}
