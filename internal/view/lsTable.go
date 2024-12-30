package view

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// LsTable ls表格界面
type LsTable struct {
	table table.Model
}

// NewLsTable 创建ls表格界面
func NewLsTable(rows []table.Row) LsTable {
	// 表头
	cols := []table.Column{
		{Title: "权限", Width: 10},
		{Title: "硬连接数", Width: 10},
		{Title: "所有者", Width: 10},
		{Title: "编辑时间", Width: 30},
		{Title: "类型", Width: 10},
		{Title: "文件名", Width: 10},
	}

	// 新建表格
	t := table.New(
		table.WithColumns(cols),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(5),
	)

	// 设置样式
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false)
	defaultStyle := lipgloss.NewStyle()
	t.SetStyles(
		table.Styles{
			Header:   borderStyle,
			Cell:     defaultStyle,
			Selected: defaultStyle,
		},
	)

	return LsTable{t}
}

func (t LsTable) Init() tea.Cmd {
	return tea.Quit
}

func (t LsTable) View() string {
	return t.table.View()
}

func (t LsTable) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	t.table, cmd = t.table.Update(msg)
	return t, cmd
}
