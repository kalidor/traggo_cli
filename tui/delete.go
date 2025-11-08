package tui

import (
	"io"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	session "github.com/kalidor/traggo_cli/session"
)

const (
	noView sessionState = iota
	yesView
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#AA0000"))
)

type deleteModel struct {
	commonModel
	deleteState sessionState
	yes         string
	no          string
	choices     []string
	taskId      int
}

func initDelete(dump io.Writer, session *session.Traggo, mainState sessionState, taskIdStr string) deleteModel {
	i, _ := strconv.Atoi(taskIdStr)
	m := deleteModel{
		commonModel: commonModel{
			dump:    dump,
			session: session,
			state:   mainState,
		},
		yes:     "YES",
		no:      "NO",
		choices: []string{"NO", "YES"},
		taskId:  i,
	}
	return m
}

func (m deleteModel) Init() tea.Cmd {
	return nil
}

func (m deleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, "deleteUpdate...")
		spew.Fdump(m.dump, msg)
		spew.Fdump(m.dump, m.state)
	}
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.deleteState == yesView {
				m.session.Delete([]int{m.taskId})
			}
			m.state = TableView

			return NewMainModel(m.dump, m.session, m.state)

		case "esc":
			m.state = TableView
			return NewMainModel(m.dump, m.session, m.state)

		case "ctrl+c", "q":
			return m, tea.Quit

		case "tab", "left", "right":
			if m.deleteState == noView {
				m.deleteState = yesView
			} else {
				m.deleteState = noView
			}
		case "y":
			m.deleteState = yesView

		case "n":
			m.deleteState = noView
		}
	}
	return m, tea.Batch(cmds...)
}

func (m deleteModel) View() string {
	var s string
	if m.deleteState == noView {
		s += lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render("NO"), modelStyle.Render("YES"))
	} else {
		s += lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render("NO"), focusedModelStyle.Render("YES"))
	}
	s += helpStyle.Render("\ny: YES • n: NO\n")
	return s
}
