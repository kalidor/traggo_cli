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

type GenericTask interface {
	GetId() int
	GetStart() time.Time
	PreparePretty(config.ColorsDef) string
}

type Error struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type CursorRequest struct {
	Offset   int `json:"offset"`
	PageSize int `json:"pageSize,omitempty"`
}

type OperationLogin struct {
	OperationName string         `json:"operationName"`
	Variables     VariablesLogin `json:"variables"`
	Query         string         `json:"query"`
}

type OperationContinue struct {
	OperationName string            `json:"operationName"`
	Variables     VariablesContinue `json:"variables,omitempty"`
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
	OldStart time.Time `json:"oldStarts,omitzero"`
	Id       int       `json:"id,omitempty"`
	Start    time.Time `json:"start,omitzero"`
	End      time.Time `json:"end,omitzero"`
	Tags     []Tag     `json:"tags,omitzero"`
	Note     string    `json:"note"` // do not omit if empty
}

type VariablesUpdateWithCursor struct {
	Start  time.Time     `json:"start"`
	End    time.Time     `json:"end"`
	Cursor CursorRequest `json:"cursor"`
}
