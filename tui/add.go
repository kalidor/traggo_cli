package tui

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kalidor/traggo_cli/config"
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
	new     bool
}

func initAdd(dump io.Writer, session *session.Traggo, mainState sessionState) addModel {
	numTags := len(session.Tags)
	var inputs []textinput.Model = make([]textinput.Model, numTags+1) // +1 for Note
	sort.Sort(config.ByPosition(session.Tags))

	for index, tag := range session.Tags {
		inputs[index] = textinput.New()
		inputs[index].Placeholder = tag.TagValueExample //"AA-1234"
		inputs[index].CharLimit = tag.CharLimit
		inputs[index].Width = tag.Width
	}

	inputs[numTags] = textinput.New()
	inputs[numTags].Placeholder = "blablabla"
	inputs[numTags].CharLimit = 100
	inputs[numTags].Width = 150
	inputs[numTags].Prompt = ""

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
		new:     false,
	}
}

func (a addModel) Init() tea.Cmd {
	return textinput.Blink
}

func (a addModel) View() string {
	helpView := a.help.View(a.keys)
	var view []string
	for index, tag := range a.session.Tags {
		view = append(view,
			fmt.Sprintf("%s: %s", inputStyle.Width(6).Render(tag.TagName), a.inputs[index].View()))
	}
	view = append(view,
		fmt.Sprintf("%s: %s", inputStyle.Width(6).Render("Note"), a.inputs[len(a.session.Tags)].View()))
	view = append(view,
		fmt.Sprintf("\n%s\n\n%s", continueStyle.Render("Continue ->"),
			helpView),
	)
	return strings.Join(view, "\n")
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
				for index, tag := range a.session.Tags {
					v := a.inputs[index].Value()
					if v == "" {
						continue
					}
					tags = append(tags, fmt.Sprintf("%s:%s", tag.TagName, v))
				}

				a.session.Start(tags, a.inputs[len(a.session.Tags)].Value())
				a.Reset()
				return NewMainModel(a.dump, a.session, a.state)

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
		if a.focused < len(a.inputs) && a.focused >= 0 {
			a.inputs[a.focused].Focus()
		}

	}
	for i := range a.inputs {
		a.inputs[i], cmds[i] = a.inputs[i].Update(msg)
	}
	return a, tea.Batch(cmds...)
}
