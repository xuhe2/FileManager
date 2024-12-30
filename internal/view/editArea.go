package view

import (
	"StarFileManager/internal/call"
	"context"
	"fmt"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
)

type EditArea struct {
	ctx      context.Context
	Filepath string
	textarea textarea.Model
}

func NewEditArea(ctx context.Context, filepath string, content string) EditArea {
	t := textarea.New()
	t.Placeholder = ""
	t.SetValue(content)
	t.Focus()
	return EditArea{ctx, filepath, t}
}
func (e EditArea) Init() tea.Cmd {
	return textarea.Blink
}

func (e EditArea) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			// 保存退出
			err := call.SaveFileContent(e.ctx, e.Filepath, e.textarea.Value())
			if err != nil {
				log.Fatalln(err)
			}
			return e, tea.Quit
		case tea.KeyCtrlC:
			// 强制退出
			return e, tea.Quit
		case tea.KeyCtrlO:
			// 保存
			err := call.SaveFileContent(e.ctx, e.Filepath, e.textarea.Value())
			if err != nil {
				log.Fatalln(err)
			}
		default:
			if !e.textarea.Focused() {
				cmd = e.textarea.Focus()
				cmds = append(cmds, cmd)
			}
		}
	case tea.WindowSizeMsg:
		e.textarea.SetWidth(msg.Width)
		e.textarea.SetHeight(msg.Height - 5)
	}

	e.textarea, cmd = e.textarea.Update(msg)
	cmds = append(cmds, cmd)
	return e, tea.Batch(cmds...)
}

func (e EditArea) View() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		e.Filepath,
		e.textarea.View(),
		"(ctrl+o to save, esc to save and quit, ctrl+c to quit but not save)",
	) + "\n\n"
}
