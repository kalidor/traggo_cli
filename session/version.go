package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Version struct {
	Name      string    `json:"name"`
	Commit    string    `json:"commit"`
	BuildDate time.Time `json:"buildDate"`
}
type RootVersion struct {
	Version Version `json:"version"`
}

type DataRootVersion struct {
	Data RootVersion `json:"data"`
}

func (t *Traggo) Version() string {
	op := OperationWithoutVariables{
		OperationName: "Version",
		Query:         "query Version {  version {    name    commit    buildDate    __typename  }}",
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
		fmt.Println("Requesting version failure")
		fmt.Println(res.Status)
	}
	var d DataRootVersion
	json.NewDecoder(res.Body).Decode(&d)

	//TODO: align output
	// Version: Name: 0.7.1
	// Commit:4aa48b385abb1728e46881964ce90a420a25f590
	// Build date:2025-04-28T15:21:13Z
	return fmt.Sprintf("Name: %s\nCommit:%s\nBuild date:%s\n", d.Data.Version.Name, d.Data.Version.Commit, d.Data.Version.BuildDate.Format(time.RFC3339))
}
