package tui

import (
	"fmt"
	"io"
	"regexp"
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
	TableView  sessionState = iota // 0
	searchView                     // 2
	periodView                     // 2
)

const (
	searchSensitive = iota
	searchInsensitive
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
	periodInput   textinput.Model
	periodString  string
	rowsOrigin    []table.Row
	lastRefreshed string
	currentTask   string
	cursor        int
	searchCase    int
}

func (m mainModel) getTasks(withComplete bool) []table.Row {
	rows := m.session.ListCurrentTasks().ToBubbleRow()
	if withComplete {
		// Add complete Tasks
		rows = append(rows, m.session.ListCompleteTasks().ToBubbleRow()...)
	}
	return rows
}

func NewMainModel(dump io.Writer, session *session.Traggo, state sessionState) (tea.Model, tea.Cmd) {
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Tags", Width: 30},
		{Title: "StartedAt", Width: 20},
		{Title: "EndedAt", Width: 20},
		{Title: "Time", Width: 10},
		{Title: "Notes", Width: 40},
	}
	rows := session.ListCurrentTasks().ToBubbleRow()
	// Add complete Tasks
	rows = append(rows, session.ListCompleteTasks().ToBubbleRow()...)
	m := mainModel{
		keys:        mainKeys,
		help:        help.New(),
		searchHelp:  help.New(),
		table:       initTable(columns, rows),
		searchInput: initSearchInput(),
		periodInput: initPeriodInput(),
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
	m.rowsOrigin = m.getTasks(true)
	m.table.SetRows(m.rowsOrigin)
	m.lastRefreshed = time.Now().Local().Format(time.DateTime)

}

// func (m *mainModel) updateDimensions(width, height int) {
// 	m.width = width
// 	m.height = height
// }

func (m *mainModel) searchInRows() {
	vSearch := m.searchInput.Value()
	if strings.Contains(vSearch, " ") {
		m.searchStrings = strings.Split(vSearch, " ")
		vSearch = ""
	}
	if len(m.searchStrings) == 0 && len(vSearch) == 0 {
		return
	}
	if len(vSearch) >= 1 {
		var sRows []table.Row
		for _, row := range m.rowsOrigin {

			tagRow := row[1]
			startDateRow := row[2]
			noteRow := row[5]
			if m.searchCase == searchInsensitive {
				tagRow = strings.ToLower(tagRow)
				noteRow = strings.ToLower(noteRow)
			}
			if len(m.searchStrings) > 0 {
				match := []bool{}
				for _, s := range m.searchStrings {
					if strings.Contains(tagRow, s) || strings.Contains(noteRow, s) || strings.Contains(startDateRow, s) {
						match = append(match, true)
					}
				}
				if strings.Contains(tagRow, vSearch) || strings.Contains(noteRow, vSearch) || strings.Contains(startDateRow, vSearch) {
					match = append(match, true)

				}
				if len(match) == len(m.searchStrings)+1 {
					sRows = append(sRows, row)
				}

			} else {
				// search in Tags and Notes
				if strings.Contains(tagRow, vSearch) || strings.Contains(noteRow, vSearch) || strings.Contains(startDateRow, vSearch) {
					sRows = append(sRows, row)
				}
			}
		}
		m.table.SetRows(sRows)
	} else {
		var sRows []table.Row

		for _, row := range m.rowsOrigin {
			tagRow := row[1]
			startDateRow := row[2]
			noteRow := row[5]
			if m.searchCase == searchInsensitive {
				tagRow = strings.ToLower(tagRow)
				noteRow = strings.ToLower(noteRow)
			}
			if len(m.searchStrings) > 0 {
				match := []bool{}

				for _, s := range m.searchStrings {
					if strings.Contains(tagRow, s) || strings.Contains(noteRow, s) || strings.Contains(startDateRow, s) {
						match = append(match, true)
					}
				}
				if len(match) == len(m.searchStrings) {
					sRows = append(sRows, row)
				}
			}
		}
		m.table.SetRows(sRows)
	}
}

func (m *mainModel) searchByPeriodInRows() {
	if m.periodString == "" {
		return
	}
	re := regexp.MustCompile(`(?P<Number>(?:-)?\d+)(?P<Type>[[:alpha:]]{1})`)
	matches := re.FindStringSubmatch(m.periodString)
	delta := func(sDate time.Time, eDate *time.Time) {}
	var startDate time.Time

	if len(matches) > 0 {
		nIndex := re.SubexpIndex("Number")
		nString := matches[nIndex]
		number, _ := strconv.Atoi(nString)
		tIndex := re.SubexpIndex("Type")

		c := matches[tIndex]
		switch c {
		case "d":
			delta = func(sDate time.Time, eDate *time.Time) {
				*eDate = sDate.AddDate(0, 0, number)
			}

		case "m":
			delta = func(sDate time.Time, eDate *time.Time) {
				*eDate = sDate.AddDate(0, number, 0)
			}

		case "w":
			delta = func(sDate time.Time, eDate *time.Time) {
				*eDate = sDate.AddDate(0, 0, number*7)
			}
		default:
			return
		}
	}
	endDate := time.Now()
	// period is negative number
	delta(endDate, &startDate)
	tasks := m.session.ListBetweenDates(startDate, endDate)
	m.table.SetRows(tasks.ToBubbleRow())
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

	case periodView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				s := m.periodInput.Value()
				if s != "" {
					m.periodString = s
				} else {
					m.state = TableView
				}
				m.periodInput.Reset()

			case "esc", "ctrl+c":
				m.state = TableView
				return m, cmd

			case "ctrl+l":
				m.periodString = ""
				m.table.SetRows(m.rowsOrigin)
				return m, cmd
			}
			m.periodInput, cmd = m.periodInput.Update(msg)

			(&m).searchByPeriodInRows()
		}
		return m, cmd
	case searchView:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				s := m.searchInput.Value()
				if s != "" {
					if !strings.Contains(s, " ") {
						m.searchStrings = append(m.searchStrings, s)
					}
				} else {
					m.state = TableView
				}
				m.searchInput.Reset()

			case "esc", "ctrl+c":
				m.state = TableView
				return m, cmd
			case "ctrl+s": // toogle search case
				if m.searchCase == searchSensitive {
					m.searchCase = searchInsensitive
					m.searchInput.Prompt = "[I]> "
				} else {
					m.searchCase = searchSensitive
					m.searchInput.Prompt = "[S]> "
				}

			case "ctrl+l":
				m.searchStrings = []string{}
				m.table.SetRows(m.rowsOrigin)
				return m, cmd
			}
			m.searchInput, cmd = m.searchInput.Update(msg)

			// TODO: ?
			// if space is in vSearch but between quote like "hello world" -> search for this word
			(&m).searchInRows()
		}
		return m, cmd

	case TableView:
		switch msg := msg.(type) {
		// case tea.WindowSizeMsg:
		// 	m.updateDimensions(msg.Width, msg.Height)
		// 	return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+w":
				// Remove last search term
				// foo bar => foo
				if len(m.searchStrings) > 0 {
					m.searchStrings = m.searchStrings[:len(m.searchStrings)-1]
					(&m).searchInRows()
				}
				// if there is no search term => display original rows
				if len(m.searchStrings) == 0 {
					m.table.SetRows(m.rowsOrigin)
				}
				return m, cmd
			case "ctrl+l":
				m.lastRefreshed = ""
				m.searchStrings = []string{}
				m.table.SetRows(m.rowsOrigin)
				return m, cmd

			case "pgup":
				m.table.MoveUp(10)
				if m.currentTask != "" {
					current_row := m.table.SelectedRow()
					if current_row == nil {
						return m, cmd
					}
					task_id, _ := strconv.Atoi(current_row[0])
					m.currentTask = m.session.SearchTask(task_id).PreparePretty(m.session.Colors)
				}
				return m, cmd
			case "pgdown":
				m.table.MoveDown(10)
				if m.currentTask != "" {
					current_row := m.table.SelectedRow()
					if current_row == nil {
						return m, cmd
					}
					task_id, _ := strconv.Atoi(current_row[0])
					m.currentTask = m.session.SearchTask(task_id).PreparePretty(m.session.Colors)
				}
				return m, cmd
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
			case "p": // period / Filter
				m.state = periodView
			case "/": // search Task / Filter
				m.state = searchView
			case "n": // add new Task
				return initEdit(m.dump, m.session, m.state, "-1").Update(msg)

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
	searchTerms := ""
	periodTerms := ""
	if len(m.searchStrings) > 0 {
		searchTerms = fmt.Sprintf("\nCurrent search: %s", strings.Join(m.searchStrings, " / "))
	}
	if m.periodString != "" {
		periodTerms = fmt.Sprintf("\nPeriod: %s", m.periodString)
		searchTerms = ""
		m.searchStrings = []string{}
	}
	switch m.state {
	case searchView:
		return baseStyle.Render(m.table.View()) + "\n" + m.searchInput.View() + searchTerms + "\n" + searchHelpView
	case periodView:
		return baseStyle.Render(m.table.View()) + "\n" + m.periodInput.View() + periodTerms + "\n" + searchHelpView
	}
	if m.currentTask != "" {
		m.currentTask = fmt.Sprintf("%s\n", m.currentTask)
	}
	if m.lastRefreshed != "" {
		m.lastRefreshed = fmt.Sprintf("Refreshed: %s\n", m.lastRefreshed)
	}

	return baseStyle.Render(m.table.View()) + "\n" + m.currentTask + searchTerms + periodTerms + "\n" + m.lastRefreshed + helpView

}

func initPeriodInput() textinput.Model {
	pti := textinput.New()
	pti.Placeholder = "Period term"
	pti.Focus()
	pti.CharLimit = 10
	pti.Width = 20
	pti.Prompt = "[Period]> "
	return pti
}

func initSearchInput() textinput.Model {
	sti := textinput.New()
	sti.Placeholder = "Search term"
	sti.Focus()
	sti.CharLimit = 156
	sti.Width = 20
	sti.Prompt = "[S]> " // default sensitive search
	return sti
}

func initTable(columns []table.Column, rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		// table.WithHeight(10),
		// table.WithWidth(100),
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
