package session

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/kalidor/traggo_cli/config"
)

// for table use only
const (
	white     = lipgloss.Color("#EEEEEE")
	lightGray = lipgloss.Color("#808080")
)

var (

	// table variables
	renderer = lipgloss.NewRenderer(os.Stdout)

	// HeaderStyle is the lipgloss style used for the table headers.
	HeaderStyle = renderer.NewStyle().Foreground(lipgloss.Color("252")).Bold(true).Align(lipgloss.Center)
	// CellStyle is the base lipgloss style used for the table rows.
	CellStyle = renderer.NewStyle().Padding(0, 1).Width(14)
	// OddRowStyle is the lipgloss style used for odd-numbered table rows.
	OddRowStyle = CellStyle.Foreground(white)
	// EvenRowStyle is the lipgloss style used for even-numbered table rows.
	EvenRowStyle = CellStyle.Foreground(lightGray)
	// BorderStyle is the lipgloss style used for the table border.
	BorderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("238"))
	baseStyle     = renderer.NewStyle().Padding(0, 1)
	SelectedStyle = baseStyle.Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
	special       = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
	checkMark     = lipgloss.NewStyle().SetString("✓").
			Foreground(special).
			PaddingRight(1).
			String()
)

type GenericTask interface {
	GetId() int
	GetStart() time.Time
	PrettyPrint(config.ColorsDef)
}

// createTimeSpan: used by Traggo.Start()
type createTimeSpanData struct {
	Data TimerTask `json:"createTimeSpan"`
}
type createTimeSpanRoot struct {
	Data   createTimeSpanData `json:"data"`
	Errors []Error            `json:"errors"`
}

type Error struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type CursorRequest struct {
	Offset   int `json:"offset"`
	PageSize int `json:"pageSize,omitempty"`
}

type Cursor struct {
	HasMore  bool `json:"hasMore"`
	StartId  int  `json:"startId,omitempty"`
	Offset   int  `json:"Offset"`
	PageSize int  `json:"pageSize,omitempty"`
}

type OperationLogin struct {
	OperationName string         `json:"operationName"`
	Variables     VariablesLogin `json:"variables"`
	Query         string         `json:"query"`
}

type OperationContinue struct {
	OperationName string            `json:"operationName"`
	Variables     VariablesContinue `json:"variables"`
	Query         string            `json:"query"`
}

type OperationBetweenDate struct {
	OperationName string                    `json:"operationName"`
	Variables     VariablesUpdateWithCursor `json:"variables"`
	Query         string                    `json:"query"`
}

type OperationUpdate struct {
	OperationName string          `json:"operationName"`
	Variables     VariablesUpdate `json:"variables"`
	Query         string          `json:"query"`
}

type OperationWithoutVariables struct {
	OperationName string `json:"operationName"`
	Query         string `json:"query"`
}

type OperationCursor struct {
	OperationName string          `json:"operationName"`
	Variables     VariablesCursor `json:"variables"`
	Query         string          `json:"query"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TimerTask struct {
	Id       int       `json:"id"`
	Start    time.Time `json:"start"`
	OldStart time.Time `json:"oldStart,omitzero"`
	Tags     []Tag     `json:"tags"`
	Note     string    `json:"note"`
	Error    string    // internal use only
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
		t.Note,
	}
}

func (t TimerTask) PrettyPrint(colors config.ColorsDef) {
	rows := [][]string{t.Export()}
	ta := table.New().
		// Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		Headers("ID", "Tags", "StartedAt", "Note").
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

			// Make the second column a little wider.
			switch col {
			case 0: // Id
				style = style.Width(5)
			case 1: // Tags
				style = style.Width(25)
				for _, c := range colors.Tags {
					if strings.Contains(rows[row][1], fmt.Sprintf("%s:%s", c.TagName, c.TagValue)) {
						return style.Width(25).Foreground(c.Color)
					}
				}
			case 2: // StartedAt
				style = style.Width(23)
			case 3: // Note
				style = style.Width(30)
			}

			return style
		}).
		Rows(rows...)
	fmt.Println(ta)
}

func (t TimerTask) String() string {
	s := fmt.Sprintf("%s [%d] \n  - start: %s\n", strings.Join(t.ExportTags(), ","), t.Id, t.Start.Format(time.DateTime))
	if len(t.Note) > 0 {
		s = fmt.Sprintf("%s  - note: %s\n", s, t.Note)
	}

	return s

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
		t.Start.Format(time.DateTime),
		t.End.Format(time.DateTime),
		duration.Round(time.Second).String(),
		strings.Join(t.ExportTags(), "\n"),
		t.Note,
	}
}

func (t TimeSpanTask) PrettyPrint(colors config.ColorsDef) {
	rows := [][]string{t.Export()}
	ta := table.New().
		// Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		Headers("ID", "StartedAt", "EndedAt", "Time", "Tags", "Note").
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

			// Change width for some column
			switch col {
			case 0: // Id
				style = style.Width(5)
			case 1, 2: // StartedAt & EndedAt
				style = style.Width(23)
			case 3: // Time
				style = style.Width(10)
			case 4: // Tags
				style = style.Width(30)
				for _, c := range colors.Tags {
					if strings.Contains(rows[row][4], fmt.Sprintf("%s:%s", c.TagName, c.TagValue)) {
						return style.Width(30).Foreground(c.Color)
					}
				}
			case 5: // Note
				style = style.Width(30)
			}

			return style
		}).
		Rows(rows...)
	fmt.Println(ta)
}

type TimeSpanTask struct {
	TimerTask
	End time.Time `json:"end,omitzero"`
}
type TimeSpanTaskList []TimeSpanTask

func (t TimeSpanTaskList) IsEmpty() bool {
	return len(t) == 0
}

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

type TimersData struct {
	Timers []TimerTask `json:"timers"`
}

func (t TimersData) IsEmpty() bool {
	return len(t.Timers) == 0
}

type TimeSpanData struct {
	Timers []TimeSpanTask `json:"timers"`
}

type TimerTasks struct {
	Data   TimersData `json:"data"`
	Errors []Error    `json:"errors"`
}

type TimeSpanTasks struct {
	Data   TimeSpanData `json:"data"`
	Errors []Error      `json:"errors"`
}

type VariablesLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
type VariablesCursor struct {
	Cursor CursorRequest `json:"cursor"`
}

type VariablesContinue struct {
	Id    int       `json:"id,omitempty"`
	Start time.Time `json:"start"`
}

type VariablesUpdate struct {
	OldStart time.Time `json:"oldStart,omitempty"`
	Id       int       `json:"id,omitempty"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Tags     []Tag     `json:"tags"`
	Note     string    `json:"note"`
}

type VariablesUpdateWithCursor struct {
	Start  time.Time     `json:"start"`
	End    time.Time     `json:"end"`
	Cursor CursorRequest `json:"cursor"`
}

func (t TimersData) PrettyPrint(colors config.ColorsDef, highlight string) {
	rows := make([][]string, len(t.Timers))

	for index, task := range t.Timers {
		rows[index] = append(rows[index], task.Export()...)
	}
	ta := table.New().
		// Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		Headers("ID", "Tags", "StartedAt", "Note").
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
				style = style.Width(25)
				for _, c := range colors.Tags {
					if strings.Contains(rows[row][1], fmt.Sprintf("%s:%s", c.TagName, c.TagValue)) {
						return style.Width(25).Foreground(c.Color)
					}
				}
			case 2: // StartedAt
				style = style.Width(23)
			case 3: // Note
				style = style.Width(30)

			}

			return style
		}).
		Rows(rows...)
	fmt.Println(ta)
}

func (t TimeSpanTaskList) PrettyPrint(colors config.ColorsDef, highlight string) {
	rows := make([][]string, len(t))
	for index, task := range t {
		rows[index] = append(rows[index], task.Export()...)
	}
	ta := table.New().
		// Border(lipgloss.ThickBorder()).
		BorderStyle(BorderStyle).
		Headers("ID", "StartedAt", "EndedAt", "Time", "Tags", "Note").
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
			case 1, 2: // StartedAt & EndedAt
				style = style.Width(23)
			case 3: // Time
				style = style.Width(10)
			case 4: // Tags
				style = style.Width(30)
				for _, c := range colors.Tags {
					if strings.Contains(rows[row][4], fmt.Sprintf("%s:%s", c.TagName, c.TagValue)) {
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
