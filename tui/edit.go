package tui

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	session "github.com/kalidor/traggo_cli/session"
	utils "github.com/kalidor/traggo_cli/utils"
)

const (
	editStartDatetime = iota
	editEndDatetime
	editTagTicket
	editTagType
	editNote
)

type editModel struct {
	commonModel
	inputs  []textinput.Model
	focused int
	task    session.GenericTask
	err     error
}

// Validator functions to ensure valid input
func datetimeValidator(s string) error {
	_, err := utils.StrToTime(s, time.DateTime)
	return err
}

func initEdit(dump io.Writer, s *session.Traggo, mainState sessionState, taskIdStr string) editModel {
	var inputs []textinput.Model = make([]textinput.Model, 5)

	taskId, _ := strconv.Atoi(taskIdStr)
	task := s.SearchTask(taskId)
	tagTicket := ""
	tagType := ""
	var tags []session.Tag
	if task.Type() == session.TypeTimerTask {
		tags = task.(session.TimerTask).Tags
	} else {
		tags = task.(session.TimeSpanTask).Tags
	}
	for _, t := range tags {
		if t.Key == "ticket" {
			tagTicket = t.Value
		}
		if t.Key == "type" {
			tagType = t.Value
		}
	}

	inputs[editStartDatetime] = textinput.New()
	inputs[editStartDatetime].Placeholder = "2025-09-12 12:00:00"
	// inputs[editStartDatetime].Focus()
	inputs[editStartDatetime].SetValue(task.GetStartString())
	inputs[editStartDatetime].CharLimit = 19
	inputs[editStartDatetime].Width = 19
	inputs[editStartDatetime].Prompt = ""
	inputs[editStartDatetime].Validate = datetimeValidator

	inputs[editEndDatetime] = textinput.New()
	inputs[editEndDatetime].Placeholder = "2025-09-12 12:00:00"
	inputs[editEndDatetime].SetValue(task.GetStopString())
	inputs[editEndDatetime].CharLimit = 19
	inputs[editEndDatetime].Width = 19
	inputs[editEndDatetime].Prompt = ""
	// No validator here since user could update current TimerTask
	// inputs[editEndDatetime].Validate = datetimeValidator

	inputs[editTagTicket] = textinput.New()
	inputs[editTagTicket].Placeholder = "HCS-1234"
	inputs[editTagTicket].CharLimit = 20
	inputs[editTagTicket].SetValue(tagTicket)
	inputs[editTagTicket].Width = 30
	inputs[editTagTicket].Prompt = ""

	inputs[editTagType] = textinput.New()
	inputs[editTagType].Placeholder = "doc/review/mdbi"
	inputs[editTagType].SetValue(tagType)
	inputs[editTagType].CharLimit = 20
	inputs[editTagType].Width = 30
	inputs[editTagType].Prompt = ""

	inputs[editNote] = textinput.New()
	inputs[editNote].Placeholder = "blablabla"
	inputs[editNote].SetValue(task.GetNote())
	inputs[editNote].CharLimit = 100
	inputs[editNote].Width = 150
	inputs[editNote].Prompt = ""

	return editModel{
		commonModel: commonModel{
			dump:    dump,
			session: s,
			state:   mainState,
		},
		task:    task,
		inputs:  inputs,
		focused: 0,
	}
}

func (e editModel) Init() tea.Cmd {
	return textinput.Blink
}

func (e editModel) View() string {
	err := ""
	if e.err != nil {
		err = fmt.Sprintf("Error: %s", e.err.Error())
		fmt.Println(err)
	}
	startErr := ""
	if e.inputs[editStartDatetime].Err != nil {
		startErr = " x"
	}
	endErr := ""
	if e.inputs[editEndDatetime].Err != nil {
		endErr = " x"
	}
	return fmt.Sprintf(
		`%s %s 
%s ->  %s

 %s: %s
 %s: %s
 %s: %s

 %s
 `,
		inputStyle.AlignHorizontal(lipgloss.Center).Width(20).Render(fmt.Sprintf("Start%s", startErr)),
		inputStyle.AlignHorizontal(lipgloss.Center).Width(25).Render(fmt.Sprintf("End%s", endErr)),
		e.inputs[editStartDatetime].View(),
		e.inputs[editEndDatetime].View(),
		inputStyle.Width(6).Render("Ticket"),
		e.inputs[editTagTicket].View(),
		inputStyle.Width(6).Render("Type"),
		e.inputs[editTagType].View(),
		inputStyle.Width(6).Render("Note"),
		e.inputs[editNote].View(),
		continueStyle.Render("Continue ->"),
	)
}

// nextInput focuses the next input field
func (e *editModel) Reset() {
	for i := range e.inputs {
		e.inputs[i].SetValue("")
	}
}

// nextInput focuses the next input field
func (e *editModel) nextInput() {
	e.focused = (e.focused + 1) % len(e.inputs)
}

// prevInput focuses the previous input field
func (e *editModel) prevInput() {
	e.focused--
	// Wrap around
	if e.focused < 0 {
		e.focused = len(e.inputs) - 1
	}
}

func (e editModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(e.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if e.focused == len(e.inputs)-1 {
				var tags []string
				ticket := e.inputs[editTagTicket].Value()
				tag := e.inputs[editTagType].Value()
				startDatetime := e.inputs[editStartDatetime].Value()
				endDatetime := e.inputs[editEndDatetime].Value()
				note := e.inputs[editNote].Value()

				if ticket == "" && tag == "" {
				} else {
					if ticket != "" {
						tags = append(tags, fmt.Sprintf("ticket:%s", ticket))
					}
					if tag != "" {
						tags = append(tags, fmt.Sprintf("type:%s", tag))
					}
					updated_task := e.task.Update(
						startDatetime,
						endDatetime,
						note,
						tags,
					)

					if endDatetime == "" {
						e.session.UpdateTimerTask(updated_task.(session.TimerTask))
					} else {
						e.session.UpdateTimeSpanTask(updated_task.(session.TimeSpanTask))
					}
					return NewMainModel(e.dump, e.session, e.state)
				}
			}
			e.nextInput()

		case tea.KeyCtrlC:
			e.Reset()
			return e, tea.Quit

		case tea.KeyShiftTab:
			e.prevInput()
		case tea.KeyTab:
			e.nextInput()
		case tea.KeyEsc:
			e.Reset()
			if e.state == TableViewCurrent {
				e.state = tableViewRefreshCurrent
			} else {
				e.state = tableViewRefreshComplete
			}

			return NewMainModel(e.dump, e.session, e.state)
		case tea.KeyCtrlL:
			e.Reset()
		}
		for i := range e.inputs {
			e.inputs[i].Blur()
		}
		e.inputs[e.focused].Focus()

	case errMsg:
		e.err = msg
	}
	for i := range e.inputs {
		e.inputs[i], cmds[i] = e.inputs[i].Update(msg)
	}
	return e, tea.Batch(cmds...)
}
