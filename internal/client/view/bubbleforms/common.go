package bubbleforms

import tea "github.com/charmbracelet/bubbletea"

const KBQuitKey = "ctrl+c"
const KBQuitDesc = "quit"
const (
	UITextMessageModelActionTemplate = "%s\n\nPress ANY key to %s."
	UITextMessageModelActionQuit     = "quit"
	UITextMessageModelActionContinue = "continue"
	UITextInputsModelSignUp          = "Sign up"
	UITextInputsModelSignIn          = "Sign in"
	UITextInputsModelAdd             = "Add"
	UITextInputsModelSaveEdit        = "Save edit"
	UITextInputsModelCancel          = "Cancel"
	UITextInputsModelLogin           = "Login"
	UITextInputsModelPassword        = "Password"
	UITextInputsModelName            = "Name"
)

type ButtonController interface {
	Handle(m tea.Model) tea.Model
}
