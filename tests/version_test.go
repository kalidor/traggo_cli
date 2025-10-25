package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kalidor/traggo_cli/config"
	session "github.com/kalidor/traggo_cli/session"
)

func TestVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type: application/json header, got: %s", r.Header.Get("Accept"))
		}
		if r.URL.Path != "/" {
			t.Errorf("Expected to request '/', got: %s", r.URL.Path)
		}

		var received session.Operation
		json.NewDecoder(r.Body).Decode(&received)

		if received.OperationName == "Version" && received.Query == "query Version {  version {    name    commit    buildDate    __typename  }}" {
			w.WriteHeader(http.StatusOK)

			w.Write([]byte(`{"data":{"version":{"name":"1.2.3", "commit":"deadbeef...", "buildDate":"2025-04-28T15:21:13Z"}}}`))
		} else {
			t.Error("Don't received expected body")
		}
	}))
	s := session.NewTraggoSession(config.NewConfigToken(server.URL, TOKEN))
	version := s.GetVersion()
	if version.Name != "1.2.3" || version.Commit != "deadbeef..." || version.BuildDate.Format(time.DateTime) != "2025-04-28 15:21:13" {
		fmt.Println(version.BuildDate.Format(time.DateTime))
		t.Error("Version data are not expected value")
	}

	defer server.Close()
}
