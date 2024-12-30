package view

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

const useHighPerformanceRenderer = false

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.BorderStyle(b)
	}()
)

type CatView struct {
	Title    string
	Content  string
	ready    bool
	Viewport viewport.Model
}

func (m CatView) Init() tea.Cmd {
	return nil
}

func (m CatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(m.footerView())
		verticalMarginHeight := headerHeight + footerHeight

		if !m.ready {
			m.Viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
			m.Viewport.YPosition = headerHeight
			m.Viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.Viewport.SetContent(m.Content)
			m.ready = true

			m.Viewport.YPosition = headerHeight + 1
		} else {
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - verticalMarginHeight
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.Viewport))
		}
	}

	m.Viewport, cmd = m.Viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m CatView) View() string {
	if !m.ready {
		return "\n  加载中..."
	}
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.Viewport.View(), m.footerView())
}

func (m CatView) headerView() string {
	title := titleStyle.Render(m.Title)
	line := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m CatView) footerView() string {
	alert := infoStyle.Border(lipgloss.NormalBorder(), false).Render("按q退出")
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.Viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.Viewport.Width-lipgloss.Width(info)-lipgloss.Width(alert)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, alert, info)
}
