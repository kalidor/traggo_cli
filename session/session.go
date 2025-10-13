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
		Colors: config.Colors,
	}
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

type TraggoUserData struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}
type TraggoUser struct {
	User *TraggoUserData `json:"user,omitempty"`
}
type TraggoCheckResponse struct {
	Data TraggoUser `json:"data"`
}

func (t *Traggo) Request(command, method string, postBody, model any) error {
	var body []byte
	body, err := json.Marshal(postBody)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(method, t.Url, bytes.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", fmt.Sprintf("traggo=%s", t.Token))

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

func (t *Traggo) Ping() error {
	op := OperationWithoutVariables{
		OperationName: "CurrentUser",
		Query:         "query CurrentUser {\n  user: currentUser {\n    name\n    id\n  }\n}\n",
	}

	var r TraggoCheckResponse
	t.Request("CurrentUser", "POST", op, &r)
	if r.Data.User == nil {
		return fmt.Errorf("successfully access traggo, but got no information about user. Check token")
	}
	return nil
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
	err = NewTraggoSession(c).Ping()
	if err != nil {
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
