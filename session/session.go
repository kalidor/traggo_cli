package session

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kalidor/traggo_cli/config"
)

type Traggo struct {
	Token  string
	Url    string
	Colors config.ColorsDef
}

func NewTraggoSession(config *config.Config) *Traggo {
	return &Traggo{
		Url:    config.Auth.Url,
		Token:  config.Auth.Token,
		Colors: config.Colors}
}

type Login struct {
	Token string `json:"token"`
}
type DataLogin struct {
	Login Login `json:"login"`
}
type TraggoAuthResponse struct {
	Data DataLogin `json:"data"`
}

func (t *Traggo) Request(command, method string, postBody, model any) error {
	var body []byte
	body, err := json.Marshal(postBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", t.Url, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", t.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode != 200 {
		c, _ := io.ReadAll(res.Body)
		fmt.Println(string(c))
		return fmt.Errorf("command '%s' failed. '%s'", command, res.Status)
	}

	json.NewDecoder(res.Body).Decode(model)

	return nil
}

func (t *Traggo) Ping() bool {
	op := OperationWithoutVariables{
		OperationName: "CurrentUser",
		Query:         "query CurrentUser {\n  user: currentUser {\n    name\n    id\n  }\n}\n",
	}
	var body []byte
	body, err := json.Marshal(op)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST", t.Url, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		log.Fatal("Ping Authentication failure")
		return false
	}
	return true
}

func RequestPermanentTokenAndTest(url, login, password string) (string, error) {
	op := OperationLogin{
		OperationName: "Login",
		Variables: VariablesLogin{
			Login:    login,
			Password: password,
		},
		Query: "mutation Login($name: String!, $pass: String!) {login(username: $name, pass: $pass, deviceName: \"test\", type: NoExpiry, cookie: false) {token user{id, name, admin, __typename}}}",
	}
	var body []byte
	body, err := json.Marshal(op)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New("authentication failure")
	}
	var d TraggoAuthResponse
	json.NewDecoder(res.Body).Decode(&d)

	// Test connectivity
	c := config.NewConfig(url, d.Data.Login.Token)
	if !NewTraggoSession(c).Ping() {
		fmt.Println("Unable to request the API")
		return "", err
	}
	return d.Data.Login.Token, nil
}

type UserSettingsData struct {
	Theme              string `json:"theme"`
	DateLocale         string `json:"dateLocale"`
	FirstDayOfTheWeek  string `json:"firstDayOfTheWeek"`
	DateTimeInputStyle string `json:"dateTimeInputStyle"`
}

type UserSettingsRoot struct {
	Data UserSettings `json:"data"`
}

type UserSettings struct {
	UserSettings UserSettingsData `json:"userSettings"`
}

func (t *Traggo) GetSettings() {
	op := OperationWithoutVariables{
		OperationName: "Settings",
		Query:         "query Settings {\n  userSettings {\n    theme\n    dateLocale\n    firstDayOfTheWeek\n    dateTimeInputStyle}\n}\n",
	}
	var body []byte
	body, err := json.Marshal(op)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", t.Url, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		log.Fatal("GetSettings failure")
	}
	var r UserSettingsRoot
	json.NewDecoder(res.Body).Decode(&r)

	fmt.Println("User settings:")
	fmt.Println("--------------")
	fmt.Printf("  - dateLocale: %s\n", r.Data.UserSettings.DateLocale)
	fmt.Printf("  - theme: %s\n", r.Data.UserSettings.Theme)
	fmt.Printf("  - firstDayOfTheWeek: %s\n", r.Data.UserSettings.FirstDayOfTheWeek)
	fmt.Printf("  - dateTimeInputStyle: %s\n", r.Data.UserSettings.DateTimeInputStyle)
}
