package tui

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type progressState struct {
	progress    progress.Model
	Current     int
	LoadingText string
	FinishText  string
}

func (p progressState) Init() tea.Cmd {
	return nil
}

func (p progressState) View()

type ProgressElement struct {
	state   *progressState
	program *tea.Program
}

func NewProgressElement(loadingText, finishText string) ProgressElement {
	return ProgressElement{
		state: &progressState{
			Current:     0,
			LoadingText: loadingText,
			FinishText:  finishText,
			progress: progress.New(progress.WithScaledGradient(MAIN_COLOR, SECOND_COLOR),
				progress.WithWidth(40)),
		},
	}
}
