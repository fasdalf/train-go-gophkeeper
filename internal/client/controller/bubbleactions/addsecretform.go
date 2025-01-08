package bubbleactions

import (
	tea "github.com/charmbracelet/bubbletea"
	ifcs "github.com/fasdalf/train-go-gophkeeper/internal/client/interfaces"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/model/entity"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/convertors"
)

var typeCheckAddSecretFormController ifcs.ButtonHandler = &AddSecretFormController{}

type AddSecretFormController struct {
	repository          ifcs.LocalRepository
	addSecretController ifcs.ButtonHandler
	cancelController    ifcs.ButtonHandler
}

func NewAddSecretFormController(
	repository ifcs.LocalRepository,
	addSecretController ifcs.ButtonHandler,
	cancelController ifcs.ButtonHandler,
) *AddSecretFormController {
	return &AddSecretFormController{
		repository:          repository,
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
