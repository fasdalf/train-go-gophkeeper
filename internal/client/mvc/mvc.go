package mvc

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/config"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/controller/bubbleactions"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/services/grpcclient"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/services/localrepository"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/services/syncroniser"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
)

var typeCheckSetTokenRemoteService bubbleactions.SetTokenRemoteService = &grpcclient.Service{}
var typeCheckRemoteService syncroniser.RemoteService = &grpcclient.Service{}
var typeCheckLocalRepository syncroniser.LocalRepository = &localrepository.LocalRepository{}
var typeCheckListSecretsLocalRepository bubbleactions.ListSecretsLocalRepository = &localrepository.LocalRepository{}
var typeCheckAddSyncroniser bubbleactions.AddSyncroniser = &syncroniser.Syncroniser{}
var typeCheckUpdateSyncroniser bubbleactions.UpdateSyncroniser = &syncroniser.Syncroniser{}

func NewModel() tea.Model {
	gService := grpcclient.NewService(config.GetConfig().GRPCAddr)
	lRepo := localrepository.NewLocalRepository()
	syncService := &syncroniser.Syncroniser{Local: lRepo, Remote: gService, SyncedUpTo: 0}
	listFormFiller := bubbleactions.NewFillListFormService(syncService, lRepo)

	cancelInputController := bubbleactions.NewCancelInputController()
	addSecretController := bubbleactions.NewAddSecretController(syncService, listFormFiller)
	addSecretFormController := bubbleactions.NewAddSecretFormController(addSecretController, cancelInputController)
	editSecretController := bubbleactions.NewEditSecretController(syncService, lRepo, listFormFiller)
	editSecretFormController := bubbleactions.NewEditSecretFormController(lRepo, editSecretController, cancelInputController)
	listModel := bubbleforms.NewListModel(addSecretFormController, editSecretFormController)
	listModelTea := tea.Model(listModel)

	signInController := bubbleactions.NewSignInController(gService, gService.SignIn, listFormFiller)
	signUpController := bubbleactions.NewSignInController(gService, gService.SignUp, listFormFiller)

	loginModel := bubbleforms.NewInputsModel(
		[]bubbleforms.InputField{
			{
				Type:  bubbleforms.InputTextPlain,
				Title: bubbleforms.UITextInputsModelLogin,
				Value: "",
			},
			{
				Type:  bubbleforms.InputTextPassword,
				Title: bubbleforms.UITextInputsModelPassword,
				Value: "",
			},
		},
		[]bubbleforms.InputButton{
			{
				Label:   bubbleforms.UITextInputsModelSignIn,
				Handler: signInController,
			},
			{
				Label:   bubbleforms.UITextInputsModelSignUp,
				Handler: signUpController,
			},
		},
		&listModelTea,
	)

	return loginModel
}
