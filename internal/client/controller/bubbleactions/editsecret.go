package bubbleactions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/convertors"
	"log/slog"
)

var typeCheckEditSecretController bubbleforms.ButtonController = &EditSecretController{}

type UpdateSyncroniser interface {
	Update(secret entity.Secret) (serverID int64, serverUpdatedAt int64, err error)
}

type GetSecretLocalRepository interface {
	GetSecret(id int) (entity.Secret, error)
}

type EditSecretController struct {
	syncroniser     UpdateSyncroniser
	repository      GetSecretLocalRepository
	listFormFiller  ListFormFiller
	secretDataValue any
}

func NewEditSecretController(
	syncroniser UpdateSyncroniser,
	repository GetSecretLocalRepository,
	listFormFiller ListFormFiller,
) *EditSecretController {
	return &EditSecretController{
		syncroniser:    syncroniser,
		repository:     repository,
		listFormFiller: listFormFiller,
	}
}

func (s *EditSecretController) setSecretDataValue(v any) {
	s.secretDataValue = v
}

func (s *EditSecretController) Handle(m tea.Model) tea.Model {
	form, ok := m.(*bubbleforms.InputsModel)
	if !ok || form == nil || len(form.Fields) == 0 {
		slog.Error("not an InputsModel or empty Fields")
		return &bubbleforms.MessageModel{Message: "not an InputsModel or empty Fields", Next: &m}
	}
	listForm, ok := (*form.ParentModel).(*bubbleforms.ListModel)
	if !ok || listForm == nil || len(listForm.List.Items()) == 0 {
		slog.Error("not a ListModel or empty list")
		return &bubbleforms.MessageModel{Message: "not a ListModel or empty list", Next: &m}
	}
	listFormTea := tea.Model(listForm)
	li := listForm.List.Index()
	secret, err := s.repository.GetSecret(li)
	if err != nil {
		slog.Error("missing secret", "error", err)
		return &bubbleforms.MessageModel{Message: "missing secret", Next: &listFormTea}
	}
	plainData, err := secret.GetData()
	if err != nil {
		slog.Error("can't decode secret", "error", err)
		return &bubbleforms.MessageModel{Message: "can't decode secret", Next: &listFormTea}
	}

	sData := entity.SecretData{}
	switch s.secretDataValue.(type) {
	case entity.Password:
		sData = convertors.FieldsToPassword(form.Fields, plainData.Metadata)
	default:
		return &bubbleforms.MessageModel{Message: "invalid secret type", Next: &listFormTea}
	}
	err = secret.SetData(sData)
	if err != nil {
		slog.Error("can't encode secret", "error", err)
		return &bubbleforms.MessageModel{Message: "can't encode secret", Next: &listFormTea}
	}
	_, _, err = s.syncroniser.Update(secret)
	if err != nil {
		slog.Error("can't send secret", "error", err)
		return &bubbleforms.MessageModel{Message: "can't send secret", Next: &m}
	}
	return s.listFormFiller.Fill(form)
}
