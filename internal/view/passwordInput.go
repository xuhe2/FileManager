package view

import (
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type PasswordInput struct {
	Username string
	input    textinput.Model
	Ctx      context.Context
	Callback func(ctx context.Context, username string, password string) error
}

func NewPasswordInput(ctx context.Context, username string, callback func(ctx context.Context, username string, password string) error) PasswordInput {
	i := textinput.New()
	i.Placeholder = ""
	i.Focus()
	return PasswordInput{
		Ctx:      ctx,
		Username: username,
		input:    i,
		Callback: callback,
	}
}

func (i PasswordInput) Init() tea.Cmd {
	return nil
}
func (i PasswordInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			// 调用回调函数
			err := i.Callback(i.Ctx, i.Username, i.input.Value())
			if err != nil {
				fmt.Println(err)
			}
			return i, tea.Quit
		}
	}
	i.input, cmd = i.input.Update(msg)
	return i, cmd
}

func (i PasswordInput) View() string {
	i.input.Focus()
	return fmt.Sprintf("Enter password:")
}
