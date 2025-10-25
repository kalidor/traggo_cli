package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
)

var currentTime = time.Date(2025, 12, 01, 00, 00, 00, 0, time.UTC)

func TestGetTokenAndTest(t *testing.T) {
	// d := t.TempDir()
	// configFileName := filepath.Join(d, "config.json")
	// create new configuration and save it to file
	variables := json.RawMessage(fmt.Sprintf(`{"Login": "%s", "Password": "%s"}`, LOGIN, PASSWORD))

	opLogin := session.Operation{
		OperationName: "Login",
		Variables:     &variables,
		Query:         "mutation Login($name: String!, $pass: String!) {login(username: $name, pass: $pass, deviceName: \"test\", type: NoExpiry, cookie: false) {token user{id, name, admin, __typename}}}",
	}

	opPing := session.Operation{
		OperationName: "CurrentUser",
		Query:         "query CurrentUser {\n  user: currentUser {\n    name\n    id\n  }\n}\n",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}

		raw, _ := io.ReadAll(r.Body)
		raw_reader := bytes.NewReader(raw)
		var dLogin session.Operation
		json.NewDecoder(raw_reader).Decode(&dLogin)

		raw_reader = bytes.NewReader(raw)
		var dPing session.Operation
		json.NewDecoder(raw_reader).Decode(&dPing)

		if dLogin.OperationName == opLogin.OperationName {
			w.WriteHeader(http.StatusOK)
			authResponse := session.TraggoAuthResponse{
				Data: session.DataLogin{
					Login: session.Login{
						Token: TOKEN,
					},
				},
			}
			var body []byte
			body, err := json.Marshal(authResponse)
			if err != nil {
				t.Fatal(err)
			}
			w.Write(body)
		} else if dPing.OperationName == opPing.OperationName {
			if r.Header.Get("Cookie") != fmt.Sprintf("traggo=%s", TOKEN) {
				t.Errorf("Expected Cookie containing token")
			}
			w.WriteHeader(http.StatusOK)
			authResponse := session.TraggoCheckResponse{
				Data: session.TraggoUser{
					User: &session.TraggoUserData{
						Name: "test_username",
						Id:   1,
					},
				},
			}
			var body []byte
			body, err := json.Marshal(authResponse)
			if err != nil {
				t.Fatal(err)
			}
			w.Write(body)

		} else {
			t.Error("Don't received expected body")
		}
	}))
	token, err := session.RequestPermanentTokenAndTest(server.URL, LOGIN, PASSWORD)
	if err != nil {
		t.Fatal(err)
	}
	if token != TOKEN {
		t.Fatalf("Expected token: '%s', got: %s", TOKEN, token)
	}
	defer server.Close()
}

func TestStart(t *testing.T) {
	session.TimeNow = func() time.Time {
		return currentTime
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}

		variables := struct {
			Start time.Time     `json:"start"`
			Tags  []session.Tag `json:"tags"`
			Note  string        `json:"note"`
		}{
			Start: session.TimeNow().UTC().Add(time.Hour * 2),
			Tags:  []session.Tag{{Key: "tag1", Value: "value1"}, {Key: "tag2", Value: "value2"}},
			Note:  "this is a note",
		}
		type expectedOperation struct {
			OperationName string `json:"operationName"`
			Variables     struct {
				Start time.Time     `json:"start"`
				Tags  []session.Tag `json:"tags"`
				Note  string        `json:"note"`
			}
			Query string `json:"query"`
		}
		expected := expectedOperation{
			OperationName: "StartTimer",
			Variables:     variables,
			Query:         "mutation StartTimer($start: Time!, $tags: [InputTimeSpanTag!], $note: String!) {\n  createTimeSpan(start: $start, tags: $tags, note: $note) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
		}

		var received expectedOperation

		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			fmt.Println(err)

		}
		if !reflect.DeepEqual(expected, received) {
			//TODO: display a diff ?
			fmt.Println(expected)
			fmt.Println(received)
			t.Errorf("POST data are not expected one.")

		}
	}))
	s := session.NewTraggoSession(config.NewConfigToken(server.URL, TOKEN))
	s.Start([]string{"tag1:value1", "tag2:value2"}, "this is a note")

	defer server.Close()
}

func TestStop(t *testing.T) {
	session.TimeNow = func() time.Time {
		return currentTime
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}

		variables := struct {
			Id  int       `json:"id"`
			End time.Time `json:"end"`
		}{
			Id:  1,
			End: session.TimeNow().UTC().Add(time.Hour * 2),
		}
		type expectedOperation struct {
			OperationName string `json:"operationName"`
			Variables     struct {
				Id  int       `json:"id"`
				End time.Time `json:"end"`
			}
			Query string `json:"query"`
		}
		expected := expectedOperation{
			OperationName: "StopTimer",
			Variables:     variables,
			Query:         "mutation StopTimer($id: Int!, $end: Time!) {\n  stopTimeSpan(id: $id, end: $end) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
		}

		var received expectedOperation

		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			fmt.Println(err)

		}
		if !reflect.DeepEqual(expected, received) {
			fmt.Println(expected)
			fmt.Println(received)
			t.Errorf("POST data are not expected one.")

		}
	}))
	s := session.NewTraggoSession(config.NewConfigToken(server.URL, TOKEN))
	s.Stop(config.ColorsDef{}, []int{1})

	defer server.Close()
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}

		variables := struct {
			Id int `json:"id"`
		}{
			Id: 1,
		}
		type expectedOperation struct {
			OperationName string `json:"operationName"`
			Variables     struct {
				Id int `json:"id"`
			}
			Query string `json:"query"`
		}
		expected := expectedOperation{
			OperationName: "RemoveTimeSpan",
			Variables:     variables,
			Query:         "mutation RemoveTimeSpan($id: Int!) {\n  removeTimeSpan(id: $id) {\n    id\n    __typename\n  }\n}\n",
		}

		var received expectedOperation

		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			fmt.Println(err)

		}
		if !reflect.DeepEqual(expected, received) {
			fmt.Println(expected)
			fmt.Println("---")
			fmt.Println(received)
			t.Errorf("POST data are not expected one.")

		}
	}))
	s := session.NewTraggoSession(config.NewConfigToken(server.URL, TOKEN))
	s.Delete([]int{1})

	defer server.Close()
}

func TestUpdateTimerStask(t *testing.T) {
	session.TimeNow = func() time.Time {
		return currentTime
	}
	task := session.TimerTask{
		Id:       123,
		OldStart: session.TimeNow(),
		Start:    session.TimeNow(),
		Tags:     []session.Tag{{Key: "tag1", Value: "value1"}},
		Note:     "this is a note",
	}

	type vUpdate struct {
		OldStart time.Time     `json:"oldStart,omitzero"`
		Id       int           `json:"id,omitempty"`
		Start    time.Time     `json:"start,omitzero"`
		End      time.Time     `json:"end,omitzero"`
		Tags     []session.Tag `json:"tags,omitzero"`
		Note     string        `json:"note"` // do not omit if empty
	}

	type opUpdate struct {
		OperationName string  `json:"operationName"`
		Variables     vUpdate `json:"variables"`
		Query         string  `json:"query"`
	}

	expected := opUpdate{
		OperationName: "UpdateTimeSpan",
		Variables: vUpdate{
			OldStart: task.OldStart,
			Id:       task.Id,
			Start:    task.Start,
			Tags:     task.Tags,
			Note:     task.Note,
		},
		Query: "mutation UpdateTimeSpan($id: Int!, $start: Time!, $tags: [InputTimeSpanTag!], $note: String!) {\n  updateTimeSpan(id: $id, start: $start, tags: $tags, note: $note) {\n    id\n    start\n    tags {\n      key\n      value\n      __typename\n    }\n   note\n    __typename\n  }\n}\n",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}

		var received opUpdate

		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			fmt.Println(err)

		}
		if !reflect.DeepEqual(expected, received) {
			fmt.Println(expected)
			fmt.Println(received)
			t.Errorf("POST data are not expected one.")

		}
	}))
	s := session.NewTraggoSession(config.NewConfigToken(server.URL, TOKEN))
	s.UpdateTimerTask(task)

	defer server.Close()
}

func TestUpdateTimeSpanTask(t *testing.T) {
	session.TimeNow = func() time.Time {
		return currentTime
	}
	task := session.TimeSpanTask{
		TimerTask: session.TimerTask{
			Id:       123,
			OldStart: session.TimeNow(),
			Start:    session.TimeNow(),
			Tags:     []session.Tag{{Key: "tag1", Value: "value1"}},
			Note:     "this is a note",
		},
		End: session.TimeNow().Add(time.Hour * 2),
	}

	type vUpdate struct {
		OldStart time.Time     `json:"oldStart,omitzero"`
		Id       int           `json:"id,omitempty"`
		Start    time.Time     `json:"start,omitzero"`
		End      time.Time     `json:"end,omitzero"`
		Tags     []session.Tag `json:"tags,omitzero"`
		Note     string        `json:"note"` // do not omit if empty
	}

	type opUpdate struct {
		OperationName string  `json:"operationName"`
		Variables     vUpdate `json:"variables"`
		Query         string  `json:"query"`
	}

	expected := opUpdate{
		OperationName: "UpdateTimeSpan",
		Variables: vUpdate{
			OldStart: task.OldStart,
			Id:       task.Id,
			Start:    task.Start,
			End:      task.End,
			Tags:     task.Tags,
			Note:     task.Note,
		},
		Query: "mutation UpdateTimeSpan($id: Int!, $start: Time!, $end: Time, $tags: [InputTimeSpanTag!], $oldStart: Time, $note: String!) {\n  updateTimeSpan(id: $id, start: $start, end: $end, tags: $tags, oldStart: $oldStart, note: $note) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}
		var received opUpdate

		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			fmt.Println(err)

		}
		if !reflect.DeepEqual(expected, received) {
			fmt.Println(expected)
			fmt.Println(received)
			t.Errorf("POST data are not expected one.")

		}
	}))
	s := session.NewTraggoSession(config.NewConfigToken(server.URL, TOKEN))
	s.UpdateTimeSpanTask(task)

	defer server.Close()
}

func TestContinue(t *testing.T) {
	session.TimeNow = func() time.Time {
		return currentTime
	}
	expected := session.Operation{
		OperationName: "Continue",
		Variables: struct {
			Id    int       `json:"id"`
			Start time.Time `json:"start"`
		}{
			Id:    123,
			Start: session.TimeNow(),
		},
		Query: "mutation Continue($id: Int!, $start: Time!) {\n  copyTimeSpan(id: $id, start: $start) {\n    id\n    start\n    __typename\n  }\n}",
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}

		var received session.Operation
		err := json.NewDecoder(r.Body).Decode(&received)
		if err != nil {
			fmt.Println(err)

		}
		if reflect.DeepEqual(expected, received) {
			fmt.Println(expected)
			fmt.Println(received)
			t.Errorf("POST data are not expected one.")
		}

	}))
	s := session.NewTraggoSession(config.NewConfigToken(server.URL, TOKEN))
	s.Continue(session.TimerTask{
		Id: 123,
	})

	defer server.Close()
}
