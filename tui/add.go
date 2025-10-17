package tui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	session "github.com/kalidor/traggo_cli/session"
)

const (
	addTagTicket = iota
	addTagType
	addNote
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type addModel struct {
	commonModel
	focused int
	help    help.Model
	inputs  []textinput.Model
	keys    addKeyMap
}

func initAdd(dump io.Writer, session *session.Traggo, mainState sessionState) addModel {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	inputs[addTagTicket] = textinput.New()
	inputs[addTagTicket].Placeholder = "AA-1234"
	// inputs[addTagTicket].Focus()
	inputs[addTagTicket].CharLimit = 20
	inputs[addTagTicket].Width = 30
	inputs[addTagTicket].Prompt = ""

	inputs[addTagType] = textinput.New()
	inputs[addTagType].Placeholder = "doc/review/mdbi"
	inputs[addTagType].CharLimit = 20
	inputs[addTagType].Width = 30
	inputs[addTagType].Prompt = ""

	inputs[addNote] = textinput.New()
	inputs[addNote].Placeholder = "blablabla"
	inputs[addNote].CharLimit = 100
	inputs[addNote].Width = 150
	inputs[addNote].Prompt = ""

	help := help.New()
	help.ShowAll = true
	return addModel{
		commonModel: commonModel{
			dump:    dump,
			session: session,
			state:   mainState,
		},
		focused: -1,
		help:    help,
		inputs:  inputs,
		keys:    addKeys,
	}
}

func (a addModel) Init() tea.Cmd {
	return textinput.Blink
}

func (a addModel) View() string {
	helpView := a.help.View(a.keys)
	return fmt.Sprintf(
		`%s: %s
%s : %s
%s: %s

%s

%s
 `,
		inputStyle.Width(6).Render("Ticket"),
		a.inputs[addTagTicket].View(),
		inputStyle.Width(6).Render("Type"),
		a.inputs[addTagType].View(),
		inputStyle.Width(6).Render("Note"),
		a.inputs[addNote].View(),
		continueStyle.Render("Continue ->"),
		helpView,
	)
}

// nextInput focuses the next input field
func (a *addModel) Reset() {
	for i := range a.inputs {
		a.inputs[i].SetValue("")
	}
}

// nextInput focuses the next input field
func (a *addModel) nextInput() {
	if a.focused == -1 {
		a.focused = 0
		return
	}
	a.focused = (a.focused + 1) % len(a.inputs)
}

// prevInput focuses the previous input field
func (a *addModel) prevInput() {
	a.focused--
	// Wrap around
	if a.focused < 0 {
		a.focused = len(a.inputs) - 1
	}
}

func (a addModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(a.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if a.focused == len(a.inputs)-1 {
				var tags []string
				ticket := a.inputs[addTagTicket].Value()
				tag := a.inputs[addTagType].Value()
				if ticket == "" && tag == "" {

				} else {
					if ticket != "" {
						tags = append(tags, fmt.Sprintf("ticket:%s", ticket))
					}
					if tag != "" {
						tags = append(tags, fmt.Sprintf("type:%s", tag))
					}

					a.session.Start(tags, a.inputs[addNote].Value())
					a.Reset()
					return NewMainModel(a.dump, a.session, a.state)

				}
			}
			a.nextInput()

		case tea.KeyCtrlC:
			return a, tea.Quit

		case tea.KeyShiftTab:
			a.prevInput()
		case tea.KeyTab:
			a.nextInput()
		case tea.KeyEsc:
			a.Reset()
			return NewMainModel(a.dump, a.session, a.state)
		case tea.KeyCtrlL:
			a.Reset()
		}
		for i := range a.inputs {
			a.inputs[i].Blur()
		}
		if a.focused != -1 {
			a.inputs[a.focused].Focus()
		}

	}
	for i := range a.inputs {
		a.inputs[i], cmds[i] = a.inputs[i].Update(msg)
	}
	return a, tea.Batch(cmds...)
}
