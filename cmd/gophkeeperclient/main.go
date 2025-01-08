package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/fasdalf/train-go-gophkeeper/internal/client/controller/bubbleactions"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/services/grpcclient"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/services/localrepository"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/services/syncroniser"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/view/bubbleforms"
	"github.com/fasdalf/train-go-gophkeeper/internal/common/printbuild"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	(&printbuild.Data{
		BuildVersion: buildVersion,
		BuildDate:    buildDate,
		BuildCommit:  buildCommit,
	}).Print()
	// TODO: ##@@ improve N
	// * Make printbuild work with --version arg only
	// * quit after output
	// * move to printbuild package
	// * Write script to fill them on build.

	// TODO: ##@@ improve N+1
	// * Write cross-platform build script for client and server.

	f, err := tea.LogToFile("debug.log", "[debug gophkeeper client]")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	gService := grpcclient.NewService("localhost:8090")
	lRepo := localrepository.NewLocalRepository()
	syncService := &syncroniser.Syncroniser{Local: lRepo, Remote: gService, SyncedUpTo: 0}
	listFormFiller := bubbleactions.NewFillListFormService(syncService, lRepo)

	cancelInputController := bubbleactions.NewCancelInputController()
	addSecretController := bubbleactions.NewAddSecretController(gService, syncService, lRepo, listFormFiller)
	addSecretFormController := bubbleactions.NewAddSecretFormController(lRepo, addSecretController, cancelInputController)
	editSecretController := bubbleactions.NewEditSecretController(gService, syncService, lRepo, listFormFiller)
	editSecretFormController := bubbleactions.NewEditSecretFormController(lRepo, editSecretController, cancelInputController)
	listModel := bubbleforms.NewListModel(addSecretFormController, editSecretFormController)
	listModelTea := tea.Model(listModel)

	signInController := bubbleactions.NewSignInController(gService, gService.SignIn, listFormFiller)
	signUpController := bubbleactions.NewSignInController(gService, gService.SignUp, listFormFiller)

	loginModel := bubbleforms.NewInputsModel(
		[]bubbleforms.InputField{
			{
				Type:  bubbleforms.InputTextPlain,
				Title: "Login",
				Value: "",
			},
			{
				Type:  bubbleforms.InputTextPassword,
				Title: "Password",
				Value: "",
			},
		},
		[]bubbleforms.InputButton{
			{
				Label:   "Sign in",
				Handler: signInController,
			},
			{
				Label:   "Sign up",
				Handler: signUpController,
			},
		},
		&listModelTea,
	)

	if _, err := tea.NewProgram(loginModel, tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
