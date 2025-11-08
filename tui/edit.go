package tui

import (
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
	utils "github.com/kalidor/traggo_cli/utils"
)

const (
	indexNote = iota + 2
	indexStartDatetime
	indexEndDatetime
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type editModel struct {
	commonModel
	inputs  []textinput.Model
	help    help.Model
	keys    editKeyMap
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
	numTags := len(s.Tags)
	var inputs []textinput.Model = make([]textinput.Model, numTags+3) // +3 for start, end and Note
	sort.Sort(config.ByPosition(s.Tags))

	taskStartString := ""
	taskStopString := ""
	taskNote := ""
	var task session.GenericTask

	if taskIdStr == "-1" {
		task = nil
		for index, tag := range s.Tags {
			inputs[index] = textinput.New()
			inputs[index].Placeholder = tag.TagValueExample //"AA-1234"
			inputs[index].CharLimit = tag.CharLimit
			inputs[index].Width = tag.Width
			inputs[index].Prompt = ""
		}
	} else {
		taskId, _ := strconv.Atoi(taskIdStr)
		task = s.SearchTask(taskId)

		taskStartString = task.GetStartString()
		taskStopString = task.GetStopString()
		taskNote = task.GetNote()

		var tags []session.Tag
		if task.Type() == session.TypeTimerTask {
			tags = task.(session.TimerTask).Tags
		} else {
			tags = task.(session.TimeSpanTask).Tags
		}

		var v string
		for index, tag := range s.Tags {
			for _, t := range tags {
				if t.Key == tag.TagName {
					v = t.Value
					break
				}
			}
			inputs[index] = textinput.New()
			inputs[index].Placeholder = tag.TagValueExample //"AA-1234"
			inputs[index].SetValue(v)
			inputs[index].CharLimit = tag.CharLimit
			inputs[index].Width = tag.Width
			inputs[index].Prompt = ""
		}
	}

	// Note
	inputs[indexNote] = textinput.New()
	inputs[indexNote].Placeholder = "blablabla"
	inputs[indexNote].SetValue(taskNote)
	inputs[indexNote].CharLimit = 100
	inputs[indexNote].Width = 150
	inputs[indexNote].Prompt = ""

	// start datetime
	inputs[indexStartDatetime] = textinput.New()
	inputs[indexStartDatetime].Placeholder = "2025-09-12 12:00:00"
	inputs[indexStartDatetime].SetValue(taskStartString)
	inputs[indexStartDatetime].CharLimit = 19
	inputs[indexStartDatetime].Width = 19
	inputs[indexStartDatetime].Prompt = ""
	inputs[indexStartDatetime].Validate = datetimeValidator

	// end datetime
	inputs[indexEndDatetime] = textinput.New()
	inputs[indexEndDatetime].Placeholder = "2025-09-12 12:00:00"
	inputs[indexEndDatetime].SetValue(taskStopString)
	inputs[indexEndDatetime].CharLimit = 19
	inputs[indexEndDatetime].Width = 19
	inputs[indexEndDatetime].Prompt = ""
	inputs[indexEndDatetime].Validate = datetimeValidator

	help := help.New()
	help.ShowAll = true
	return editModel{
		commonModel: commonModel{
			dump:    dump,
			session: s,
			state:   mainState,
		},
		task:    task,
		inputs:  inputs,
		help:    help,
		keys:    editKeys,
		focused: -1,
	}
}

func (e editModel) Init() tea.Cmd {
	return textinput.Blink
}

func (e editModel) View() string {
	helpView := e.help.View(e.keys)

	err := ""
	if e.err != nil {
		err = fmt.Sprintf("Error: %s", e.err.Error())
		fmt.Println(err)
	}
	startErr := ""
	if e.inputs[indexStartDatetime].Err != nil {
		startErr = " x"
	}
	endErr := ""
	if e.inputs[indexEndDatetime].Err != nil {
		endErr = " x"
	}

	var view []string

	for index, tag := range e.session.Tags {
		view = append(view,
			fmt.Sprintf("%s: %s", inputStyle.Width(6).Render(strings.ToLower(tag.TagName)), e.inputs[index].View()))
	}
	// Note
	view = append(view,
		fmt.Sprintf("%s: %s", inputStyle.Width(6).Render("Note"), e.inputs[indexNote].View()))

	// start datetime
	view = append(view,
		fmt.Sprintf("%s: %s", inputStyle.Width(8).Render(fmt.Sprintf("Start%s", startErr)), e.inputs[indexStartDatetime].View()))

	// end datetime
	view = append(view,
		fmt.Sprintf("%s: %s", inputStyle.Width(8).Render(fmt.Sprintf("End%s", endErr)), e.inputs[indexEndDatetime].View()))
	view = append(view,
		fmt.Sprintf("\n%s\n\n%s", continueStyle.Render("Continue ->"), helpView),
	)
	return strings.Join(view, "\n")
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
				for index, tag := range e.session.Tags {
					v := e.inputs[index].Value()
					if v == "" {
						continue
					}
					tags = append(tags, fmt.Sprintf("%s:%s", strings.ToLower(tag.TagName), v))
				}
				// it's a new task
				if e.task == nil {
					e.session.Start(tags, e.inputs[len(e.session.Tags)].Value())
					e.Reset()
					return NewMainModel(e.dump, e.session, e.state)
				}

				// otherwise it's task update

				if len(tags) != 0 {
					endDatetime := e.inputs[indexEndDatetime].Value()
					updated_task := e.task.Update(
						e.inputs[indexStartDatetime].Value(),
						endDatetime,
						e.inputs[indexNote].Value(),
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
			return e, tea.Quit

		case tea.KeyShiftTab, tea.KeyUp:
			e.prevInput()
		case tea.KeyTab, tea.KeyDown:
			e.nextInput()
		case tea.KeyEsc:
			e.Reset()
			e.state = TableView

			return NewMainModel(e.dump, e.session, e.state)
		case tea.KeyCtrlL:
			e.Reset()
		}
		for i := range e.inputs {
			e.inputs[i].Blur()
		}

		if e.focused < len(e.inputs) && e.focused >= 0 {
			e.inputs[e.focused].Focus()
		}

	case errMsg:
		e.err = msg
	}
	for i := range e.inputs {
		e.inputs[i], cmds[i] = e.inputs[i].Update(msg)
	}
	return e, tea.Batch(cmds...)
}
