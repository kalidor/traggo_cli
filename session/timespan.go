package session

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/kalidor/traggo_cli/config"
)

type Cursor struct {
	HasMore  bool `json:"hasMore"`
	StartId  int  `json:"startId,omitempty"`
	Offset   int  `json:"Offset"`
	PageSize int  `json:"pageSize,omitempty"`
}

type TimeSpanTask struct {
	TimerTask
	End time.Time `json:"end,omitzero"`
}

type TimeSpanTaskList []TimeSpanTask

type TimeSpans struct {
	TimeSpans TimeSpanTaskList `json:"timeSpans"`
	Cursor    Cursor           `json:"cursor"`
}
type TimeSpansData struct {
	TimeSpans TimeSpans `json:"timeSpans"`
}
type TimeSpanRoot struct {
	Data TimeSpansData `json:"data"`
}

type TimeSpanData struct {
	Timers []TimeSpanTask `json:"timers"`
}

type TimeSpanTasks struct {
	Data   TimeSpanData `json:"data"`
	Errors []Error      `json:"errors"`
}

func (t TimeSpanTask) GetId() int {
	return t.Id
}

func (t TimeSpanTask) GetStart() time.Time {
	return t.Start
}

func (t TimeSpanTask) String() string {
	duration := t.End.Sub(t.Start)
	s := fmt.Sprintf("%s [%d] \n  - start: %s\n  - started from now: %s\n  - end: %s\n", strings.Join(t.ExportTags(), ","), t.Id, t.Start.Format(time.DateTime), duration.Round(time.Second).String(), t.End.Format(time.DateTime))
	if len(t.Note) > 0 {
		s = fmt.Sprintf("%s  - note: %s\n", s, t.Note)
	}
	return s
}

func (t TimeSpanTask) Export() []string {
	duration := t.End.Sub(t.Start)

	return []string{
		fmt.Sprintf("%d", t.Id),
		strings.Join(t.ExportTags(), "\n"),
		t.Start.Format(time.DateTime),
		t.End.Format(time.DateTime),
		duration.Round(time.Second).String(),
		t.Note,
	}
}

func (t TimeSpanTask) PrettyPrint(colors config.ColorsDef) {
	var l TimeSpanTaskList
	l = append(l, t)
	l.PrettyPrint(colors, "")
}

func (t TimeSpanTaskList) IsEmpty() bool {
	return len(t) == 0
}

func (t TimeSpanTaskList) PrettyPrint(colors config.ColorsDef, highlight string) {
	rows := make([][]string, len(t))
	for index, task := range t {
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
			if highlight != "" && (strings.Contains(rows[row][3], highlight) || strings.Contains(rows[row][4], highlight)) {
				return SelectedStyle
			}

			// Change width for some column
			switch col {
			case 0: // Id
				style = style.Width(5)
			case 2, 3: // StartedAt & EndedAt
				style = style.Width(23)
			case 4: // Time
				style = style.Width(10)
			case 1: // Tags
				style = style.Width(30)
				for _, c := range colors.Tags {
					if strings.Contains(rows[row][1], fmt.Sprintf("%s:%s", c.TagName, c.TagValue)) {
						return style.Width(30).Foreground(c.Color)
					}
				}
			case 5: // Note
				style = style.Width(30)
			}

			return style
		}).
		Rows(rows...)

	// finally display it
	fmt.Println(ta)
}
