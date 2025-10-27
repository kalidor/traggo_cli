package tui

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/davecgh/go-spew/spew"
	session "github.com/kalidor/traggo_cli/session"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type sessionState int

const (
	TableViewCurrent         sessionState = iota // 0
	TableViewComplete                            // 1
	searchView                                   // 2
	editView                                     // 3
	addView                                      // 4
	refreshView                                  // 5
	tableViewRefreshCurrent                      // 6
	tableViewRefreshComplete                     // 7

)

type errMsg struct{ error }

type commonModel struct {
	dump    io.Writer
	session *session.Traggo
	state   sessionState
}

type mainModel struct {
	commonModel
	keys          mainKeyMap
	help          help.Model
	searchHelp    help.Model
	table         table.Model
	searchInput   textinput.Model
	searchStrings []string
	rowsOrigin    []table.Row
	lastRefreshed string
	currentTask   string
	cursor        int
	previousState sessionState
}

func NewMainModel(dump io.Writer, session *session.Traggo, state sessionState) (tea.Model, tea.Cmd) {
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Tags", Width: 35},
		{Title: "StartedAt", Width: 20},
		{Title: "EndedAt", Width: 20},
	}
	rows := session.ListCurrentTasks().ToBubbleRow()
	m := mainModel{
		keys:        mainKeys,
		help:        help.New(),
		searchHelp:  help.New(),
		table:       initTable(columns, rows),
		searchInput: initSearchInput(),
		rowsOrigin:  rows,
		commonModel: commonModel{
			dump:    dump,
			session: session,
			state:   state,
		},
	}
	return m, func() tea.Msg { return errMsg{nil} }
}

func (m mainModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m *mainModel) Refresh() {
	switch m.state {
	case TableViewCurrent:
		m.rowsOrigin = m.session.ListCurrentTasks().ToBubbleRow()
	case TableViewComplete:
		m.rowsOrigin = m.session.ListCompleteTasks().ToBubbleRow()
	}
	m.table.SetRows(m.rowsOrigin)
	m.lastRefreshed = time.Now().Format(time.DateTime)

}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.dump != nil {
		spew.Fdump(m.dump, msg)
		spew.Fdump(m.dump, m.state)
		spew.Fdump(m.dump, m.cursor)
		spew.Fdump(m.dump, m.table.Cursor())
	}
	var cmd tea.Cmd

	switch m.state {

	case tableViewRefreshCurrent, tableViewRefreshComplete:
		if m.state == tableViewRefreshCurrent {
			m.state = TableViewCurrent
		} else {
			m.state = TableViewComplete
		}
		m.Refresh()
		return m, cmd
	case searchView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				s := m.searchInput.Value()
				if s != "" {
					m.searchStrings = append(m.searchStrings, s)
				} else {
					m.state = m.previousState
				}
				m.searchInput.Reset()

			case "esc":
				m.state = m.previousState

			case "ctrl+l":
				m.searchStrings = []string{}
				m.table.SetRows(m.rowsOrigin)
			}
			m.searchInput, cmd = m.searchInput.Update(msg)
			vSearch := m.searchInput.Value()
			// TODO:
			// if space is in vSearch -> search for both word separately
			// if space is in vSearch but between quote like "hello world" -> search for this word
			if len(vSearch) >= 1 {
				var sRows []table.Row
				for _, row := range m.rowsOrigin {
					if len(m.searchStrings) > 0 {
						for _, s := range m.searchStrings {
							if strings.Contains(row[1], s) && strings.Contains(row[1], vSearch) {
								sRows = append(sRows, row)
							}
						}
					} else {
						if strings.Contains(row[1], vSearch) {
							sRows = append(sRows, row)
						}
					}
				}
				m.table.SetRows(sRows)
			}
		}
		return m, cmd

	case TableViewComplete, TableViewCurrent:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "tab":
				if m.state == TableViewCurrent {
					m.state = TableViewComplete
				} else {
					m.state = TableViewCurrent
				}
				m.Refresh()

			case "ctrl+l":
				m.lastRefreshed = ""

			case "up":
				m.table.MoveUp(1)
				if m.currentTask != "" {
					current_row := m.table.SelectedRow()
					if current_row == nil {
						return m, cmd
					}
					task_id, _ := strconv.Atoi(current_row[0])
					m.currentTask = m.session.SearchTask(task_id).PreparePretty(m.session.Colors)
				}
				return m, cmd
			case "down":
				m.table.MoveDown(1)
				if m.currentTask != "" {
					current_row := m.table.SelectedRow()
					if current_row == nil {
						return m, cmd
					}
					task_id, _ := strconv.Atoi(current_row[0])
					m.currentTask = m.session.SearchTask(task_id).PreparePretty(m.session.Colors)
				}
				return m, cmd
			case "q", "ctrl+c", "esc":
				if m.currentTask != "" {
					m.currentTask = ""
				} else {
					return m, tea.Quit
				}
			case "/": // search Task / Filter
				m.previousState = m.state
				m.state = searchView
			case "n": // add new Task
				return initAdd(m.dump, m.session, m.state).Update(msg)

			case "d": // delete Task
				current_row := m.table.SelectedRow()
				if current_row == nil {
					return m, cmd
				}
				return initDelete(m.dump, m.session, m.state, current_row[0]).Update(msg)

			case "e", "u": // edit & update
				current_row := m.table.SelectedRow()
				if current_row == nil {
					return m, cmd
				}
				return initEdit(m.dump, m.session, m.state, current_row[0]).Update(msg)
			case "c": // continue
				current_row := m.table.SelectedRow()
				if current_row == nil {
					return m, cmd
				}
				taskId, _ := strconv.Atoi(current_row[0])
				m.session.Continue(m.session.SearchTask(taskId))
				m.Refresh()
			case "s": // stop
				current_row := m.table.SelectedRow()
				if current_row == nil {
					return m, cmd
				}
				taskId, _ := strconv.Atoi(current_row[0])
				m.session.Stop(m.session.Colors, []int{taskId})
				m.Refresh()
			case "r": // refresh
				m.searchStrings = []string{}
				m.Refresh()

			case "?":
				m.help.ShowAll = !m.help.ShowAll
			case "enter":
				m.lastRefreshed = ""

				if m.currentTask != "" {
					m.currentTask = ""
				} else {
					current_row := m.table.SelectedRow()
					if current_row == nil {
						return m, cmd
					}
					task_id, _ := strconv.Atoi(current_row[0])
					m.currentTask = m.session.SearchTask(task_id).PreparePretty(m.session.Colors)
				}
			}
		}

		return m, cmd

	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m mainModel) View() string {
	helpView := m.help.View(m.keys)
	searchHelpView := m.searchHelp.View(searchKeys)
	seachTerms := ""
	if len(m.searchStrings) > 0 {
		seachTerms = fmt.Sprintf("\nCurrent search: %s", strings.Join(m.searchStrings, " / "))
	}
	switch m.state {
	case searchView:
		return baseStyle.Render(m.table.View()) + "\n" + m.searchInput.View() + seachTerms + "\n" + searchHelpView
	}
	if m.currentTask != "" {
		m.currentTask = fmt.Sprintf("%s\n", m.currentTask)
	}
	if m.lastRefreshed != "" {
		m.lastRefreshed = fmt.Sprintf("Refreshed: %s\n", m.lastRefreshed)
	}

	return baseStyle.Render(m.table.View()) + "\n" + m.currentTask + seachTerms + "\n" + m.lastRefreshed + helpView

}

func initSearchInput() textinput.Model {
	sti := textinput.New()
	sti.Placeholder = "Pikachu"
	sti.Focus()
	sti.CharLimit = 156
	sti.Width = 20
	return sti
}

func initTable(columns []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(10),
		table.WithWidth(90),
	)

	style := table.DefaultStyles()
	style.Header = style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	style.Selected = style.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(style)
	return t
}
