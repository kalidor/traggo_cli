package session

import (
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
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

type taskType int

const (
	TypeTimerTask taskType = iota
	TypeTimeSpanTask
)

type GenericTask interface {
	GetId() int
	GetNote() string
	GetStart() time.Time
	GetStartString() string
	GetStopString() string
	PreparePretty(config.ColorsDef) string
	Type() taskType
	Update(start, stop, note string, tags []string) GenericTask
}

type Error struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type CursorRequest struct {
	Offset   int `json:"offset"`
	PageSize int `json:"pageSize,omitempty"`
}

type Operation struct {
	OperationName string `json:"operationName"`
	Variables     any    `json:"variables,omitempty"`
	Query         string `json:"query"`
}

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
