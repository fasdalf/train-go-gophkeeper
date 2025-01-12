package bubbleactions

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
	"log/slog"
)

var typeCheckSignInController bubbleforms.ButtonController = &SignInController{}

type SetTokenRemoteService interface {
	SetToken(token string)
}

type getTokenFunc func(login, pass string) (string, error)
type SignInController struct {
	service        SetTokenRemoteService
	getToken       getTokenFunc
	listFormFiller ListFormFiller
}

func NewSignInController(
	service SetTokenRemoteService,
	getToken getTokenFunc,
	listFormFiller ListFormFiller,
) *SignInController {
	return &SignInController{
		service:        service,
		getToken:       getToken,
		listFormFiller: listFormFiller,
	}
}

func (s *SignInController) Handle(m tea.Model) tea.Model {
	form, ok := m.(*bubbleforms.InputsModel)
	if !ok || form == nil || len(form.Fields) < 2 {
		slog.Error("could not cast model to InputsModel")
		return &bubbleforms.MessageModel{Message: "could not cast model to InputsModel", Next: &m}
	}
	login := form.Fields[0].Value
	pass := form.Fields[1].Value
	token, err := s.getToken(login, pass)
	if err != nil {
		slog.Error("could not sign in user", "err", err)
		return &bubbleforms.MessageModel{Message: err.Error(), Next: &m}
	}
	s.service.SetToken(token)
	return s.listFormFiller.Fill(form)
}
