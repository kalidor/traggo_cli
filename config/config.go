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

// Config contains all configuration related information
type Config struct {
	Auth   Auth      `json:"auth"`   // use for authentication
	Colors ColorsDef `json:"colors"` // use for user experience, to colorize output for matching tags
}

func NewConfig(url, token string) *Config {
	return &Config{
		Auth{Url: url, Token: token},
		ColorsDef{},
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

func NewConfigToken(url string, token string) *Config {
	return &Config{
		Auth{Url: url, Token: token},
		ColorsDef{},
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
