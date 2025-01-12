package main

import (
	"fmt"
	"github.com/fasdalf/train-go-gophkeeper/internal/client/mvc"
	"os"

	tea "github.com/charmbracelet/bubbletea"

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

	if _, err := tea.NewProgram(mvc.NewModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
