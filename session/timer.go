package session

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/kalidor/traggo_cli/config"
)

type TimerTask struct {
	Id       int       `json:"id"`
	Start    time.Time `json:"start"`
	OldStart time.Time `json:"oldStart,omitzero"`
	Tags     []Tag     `json:"tags"`
	Note     string    `json:"note"`
	Error    string    // internal use only
}

type TimersData struct {
	Timers []TimerTask `json:"timers"`
}

func (t TimerTask) ExportTags() []string {
	var r []string
	for _, tags := range t.Tags {
		r = append(r, fmt.Sprintf("%s:%s", tags.Key, tags.Value))
	}
	return r
}

func (t TimerTask) GetId() int {
	return t.Id
}

func (t TimerTask) GetStart() time.Time {
	return t.Start
}

func (t TimerTask) Export() []string {
	return []string{
		fmt.Sprintf("%d", t.Id),
		strings.Join(t.ExportTags(), "\n"),
		t.Start.Format(time.DateTime),
		"-", "-",
		t.Note,
	}
}

func (t TimerTask) PrettyPrint(colors config.ColorsDef) {
	var l TimersData
	l.Timers = append(l.Timers, t)
	l.PrettyPrint(colors, "")
}

func (t TimerTask) String() string {
	s := fmt.Sprintf("%s [%d] \n  - start: %s\n", strings.Join(t.ExportTags(), ","), t.Id, t.Start.Format(time.DateTime))
	if len(t.Note) > 0 {
		s = fmt.Sprintf("%s  - note: %s\n", s, t.Note)
	}

	return s
}

func (t TimersData) IsEmpty() bool {
	return len(t.Timers) == 0
}

type TimerTasks struct {
	Data   TimersData `json:"data"`
	Errors []Error    `json:"errors"`
}

func (t TimersData) PrettyPrint(colors config.ColorsDef, highlight string) {
	rows := make([][]string, len(t.Timers))

	for index, task := range t.Timers {
		rows[index] = append(rows[index], task.Export()...)
	}
	ta := table.New().
		// Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		Headers("ID", "Tags", "StartedAt", "EndedAt", "Time", "Note").
		StyleFunc(func(row, col int) lipgloss.Style {
			var style lipgloss.Style

			switch {
			case row == table.HeaderRow:
				return baseStyle.Foreground(colors.Table.HeaderStyle).Bold(true)
			case row%2 == 0:
				style = CellStyle.Foreground(colors.Table.EvenStyle)
			default:
				style = CellStyle.Foreground(colors.Table.OddStyle)
			}
			if highlight != "" && (strings.Contains(rows[row][1], highlight) || strings.Contains(rows[row][3], highlight)) {
				return SelectedStyle
			}

			// Make the second column a little wider.
			switch col {
			case 0: //id
				style = style.Width(5)
			case 1: // Tags
				style = style.Width(30)
				for _, c := range colors.Tags {
					if strings.Contains(rows[row][1], fmt.Sprintf("%s:%s", c.TagName, c.TagValue)) {
						return style.Width(25).Foreground(c.Color)
					}
				}
			case 2, 3: // StartedAt
				style = style.Width(23)
			case 4: // Time
				style = style.Width(10)
			case 5: // Note
				style = style.Width(30)

			}

			return style
		}).
		Rows(rows...)
	fmt.Println(ta)
}
