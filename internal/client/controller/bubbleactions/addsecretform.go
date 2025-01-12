package bubbleactions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/convertors"
)

var typeCheckAddSecretFormController bubbleforms.ButtonController = &AddSecretFormController{}

type AddSecretFormController struct {
	addSecretController bubbleforms.ButtonController
	cancelController    bubbleforms.ButtonController
}

func NewAddSecretFormController(
	addSecretController bubbleforms.ButtonController,
	cancelController bubbleforms.ButtonController,
) *AddSecretFormController {
	return &AddSecretFormController{
		addSecretController: addSecretController,
		cancelController:    cancelController,
	}
}

func (s *AddSecretFormController) Handle(m tea.Model) tea.Model {
	plainData := entity.SecretData{Value: entity.Password{}}

	fields := []bubbleforms.InputField{}
	buttons := []bubbleforms.InputButton{
		{
			Label:   "Add",
			Handler: s.addSecretController,
		},
		{
			Label:   "Cancel",
			Handler: s.cancelController,
		},
	}
	switch plainData.Value.(type) {
	case entity.Password:
		fields = convertors.PasswordToFields(plainData)
	default:
		return &bubbleforms.MessageModel{Message: "invalid secret type", Next: &m}
	}

	ac := s.addSecretController.(*AddSecretController)
	if ac != nil {
		ac.setSecretDataValue(plainData.Value)
	}

	return bubbleforms.NewInputsModel(
		fields,
		buttons,
		&m,
	)
}
