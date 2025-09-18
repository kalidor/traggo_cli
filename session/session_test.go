package session

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	URL      = "https://my-traggo.io"
	LOGIN    = "username"
	PASSWORD = "s3cr3tP4ssw0rd"
	TOKEN    = "myfreshlynewtoken"
)

func TestGetTokenAndTest(t *testing.T) {
	// d := t.TempDir()
	// configFileName := filepath.Join(d, "config.json")
	// create new configuration and save it to file
	opLogin := OperationLogin{
		OperationName: "Login",
		Variables: VariablesLogin{
			Login:    LOGIN,
			Password: PASSWORD,
		},
		Query: "mutation Login($name: String!, $pass: String!) {login(username: $name, pass: $pass, deviceName: \"test\", type: NoExpiry, cookie: false) {token user{id, name, admin, __typename}}}",
	}

	opPing := OperationWithoutVariables{
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
		var dLogin OperationLogin
		json.NewDecoder(raw_reader).Decode(&dLogin)

		raw_reader = bytes.NewReader(raw)
		var dPing OperationWithoutVariables
		json.NewDecoder(raw_reader).Decode(&dPing)

		if dLogin == opLogin {
			w.WriteHeader(http.StatusOK)
			authResponse := TraggoAuthResponse{
				Data: DataLogin{
					Login: Login{
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
		} else if dPing == opPing {
			w.WriteHeader(http.StatusOK)

		} else {
			t.Error("Don't received expected body")
		}
	}))
	token, err := RequestPermanentTokenAndTest(server.URL, LOGIN, PASSWORD)
	if err != nil {
		t.Fatal(err)
	}
	if token != TOKEN {
		t.Fatalf("Expected token: '%s', got: %s", TOKEN, token)
	}
	defer server.Close()
}
