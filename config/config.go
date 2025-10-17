package config

import (
	"encoding/json"
	"log"
	"os"
	"path"

	"github.com/charmbracelet/lipgloss"
)

type Auth struct {
	Url   string `json:"url"`   // endpoint URL
	Token string `json:"token"` // Token is retrieved the very first time then store in configuration. All future calls will use it
}

// ColorTagDef will allow colorize matching tagName and tagValue
type ColorTagDef struct {
	TagName  string         `json:"tagName"`
	TagValue string         `json:"tagValue"`
	Color    lipgloss.Color `json:"color"`
}

type ColorTagDefs []ColorTagDef

type ColorTableDef struct {
	EvenStyle   lipgloss.Color `json:"even"`
	HeaderStyle lipgloss.Color `json:"header"`
	OddStyle    lipgloss.Color `json:"odd"`
}

type ColorsDef struct {
	Tags  ColorTagDefs  `json:"tags"`
	Table ColorTableDef `json:"table"`
}

type TagDef struct {
	TagName         string `json:"tagName"`
	TagValueExample string `json:"tagValueExample"`
	Position        int    `json:"position"`
	CharLimit       int    `json:"charLimit"`
	Width           int    `json:"width"`
}
type TagsDef []TagDef

// ByPosition implements sort.Interface for []TagDef based on
// the Position field.
type ByPosition []TagDef

func (a ByPosition) Len() int           { return len(a) }
func (a ByPosition) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPosition) Less(i, j int) bool { return a[i].Position < a[j].Position }

// Config contains all configuration related information
type Config struct {
	Auth   Auth      `json:"auth"`   // use for authentication
	Colors ColorsDef `json:"colors"` // use for user experience, to colorize output for matching tags
	Tags   TagsDef   `json:"tags"`   // use for user experience, to specify how many tags should be proposed in "live" mode
}

func NewConfig(url, token string) *Config {
	return &Config{
		Auth{Url: url, Token: token},
		ColorsDef{},
		TagsDef{},
	}
}

func LoadConfig(configPath string) *Config {
	var c Config
	d, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(d, &c)
	if err != nil {
		panic(err)
	}
	return &c
}

// TODO: to remove, because not used
func NewConfigToken(url string, token string) *Config {
	return &Config{
		Auth{Url: url, Token: token},
		ColorsDef{},
		TagsDef{},
	}
}

func (c *Config) Save(configPath string) error {
	configDir := path.Dir(configPath)

	err := os.MkdirAll(configDir, 0770)
	if err != nil {
		return err
	}

	f, err := os.Create(configPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	f.Write(data)
	return nil
}
