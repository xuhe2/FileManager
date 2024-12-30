package model

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

// LsTable ls表格界面
type LsTable struct {
	Table table.Model
}

func (t LsTable) Init() tea.Cmd {
	return tea.Quit
}

func (t LsTable) View() string {
	return t.Table.View()
}

func (t LsTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	t.Table, cmd = t.Table.Update(msg)
	return t, cmd
}
