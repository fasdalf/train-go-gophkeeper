package bubbleactions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/convertors"
	"log/slog"
)

var typeCheckEditSecretFormController bubbleforms.ButtonController = &EditSecretFormController{}

type EditSecretFormController struct {
	repository           GetSecretLocalRepository
	editSecretController bubbleforms.ButtonController
	cancelController     bubbleforms.ButtonController
}

func NewEditSecretFormController(
	repository GetSecretLocalRepository,
	editSecretController bubbleforms.ButtonController,
	cancelController bubbleforms.ButtonController,
) *EditSecretFormController {
	return &EditSecretFormController{
		repository:           repository,
		editSecretController: editSecretController,
		cancelController:     cancelController,
	}
}

func (s *EditSecretFormController) Handle(m tea.Model) tea.Model {
	form, ok := m.(*bubbleforms.ListModel)
	if !ok || form == nil || len(form.List.Items()) == 0 {
		slog.Error("not a ListModel or empty list", "ok", ok, "model", form)
		return &bubbleforms.MessageModel{Message: "not a ListModel or empty list", Next: &m}
	}
	li := form.List.Index()
	secret, err := s.repository.GetSecret(li)
	if err != nil {
		slog.Error("missing secret", "error", err)
		return &bubbleforms.MessageModel{Message: "missing secret", Next: &m}
	}
	plainData, err := secret.GetData()
	if err != nil {
		slog.Error("can't decode secret", "error", err)
		return &bubbleforms.MessageModel{Message: "can't decode secret", Next: &m}
	}

	fields := []bubbleforms.InputField{}
	buttons := []bubbleforms.InputButton{
		{
			Label:   bubbleforms.UITextInputsModelSaveEdit,
			Handler: s.editSecretController,
		},
		{
			Label:   bubbleforms.UITextInputsModelCancel,
			Handler: s.cancelController,
		},
	}
	switch plainData.Value.(type) {
	case entity.Password:
		fields = convertors.PasswordToFields(plainData)
	default:
		return &bubbleforms.MessageModel{Message: "invalid secret type", Next: &m}
	}

	ac := s.editSecretController.(*EditSecretController)
	if ac != nil {
		ac.setSecretDataValue(plainData.Value)
	}

	return bubbleforms.NewInputsModel(
		fields,
		buttons,
		&m,
	)
}
