package session

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kalidor/traggo_cli/config"
)

type createTimeSpanData struct {
	Data TimerTask `json:"createTimeSpan"`
}
type createTimeSpanRoot struct {
	Data   createTimeSpanData `json:"data"`
	Errors []Error            `json:"errors"`
}

func (t *Traggo) Start(tags []string, note string) {
	var splitedSlice []Tag
	for _, tag := range tags {
		if strings.Contains(tag, ":") {
			s := strings.SplitN(tag, ":", 2)
			splitedSlice = append(splitedSlice, Tag{Key: s[0], Value: s[1]})
		}
	}

	op := OperationUpdate{
		OperationName: "StartTimer",
		Variables: VariablesUpdate{
			Start: time.Now(),
			Tags:  splitedSlice,
			Note:  note,
		},
		Query: "mutation StartTimer($start: Time!, $tags: [InputTimeSpanTag!], $note: String!) {\n  createTimeSpan(start: $start, tags: $tags, note: $note) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
	}

	// Parse http.Response Boby as JSON and display it
	var d createTimeSpanRoot
	err := t.Request(
		"Start",
		"POST",
		op,
		&d,
	)
	if err != nil {
		log.Fatal(err)
	}
	d.Data.Data.PrettyPrint(t.Colors)

}

func (t *Traggo) Stop(colors config.ColorsDef, ids []int) {
	op := OperationUpdate{
		OperationName: "StopTimer",
		Variables: VariablesUpdate{
			Id:  0,
			End: time.Now(),
		},
		Query: "mutation StopTimer($id: Int!, $end: Time!) {\n  stopTimeSpan(id: $id, end: $end) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
	}

	for _, id := range ids {
		fmt.Printf("Stopping: %d \n", id)
		op.Variables.Id = id
		var d TimeSpanTask
		err := t.Request(
			"ListBetweenDates",
			"POST",
			op,
			&d,
		)
		if err != nil {
			log.Fatal(err)
		}
		d.PrettyPrint(colors)
	}
}

func (t *Traggo) Delete(ids []int) {
	op := OperationUpdate{
		OperationName: "RemoveTimeSpan",
		Variables: VariablesUpdate{
			Id: 0,
		},
		Query: "mutation RemoveTimeSpan($id: Int!) {\n  removeTimeSpan(id: $id) {\n    id\n    __typename\n  }\n}\n",
	}
	for _, id := range ids {
		op.Variables.Id = id
		body, err := json.Marshal(op)
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
			log.Fatal("Deleting task failure")
		}
	}
}

func (t *Traggo) UpdateTimerTask(task TimerTask) {
	op := OperationUpdate{
		OperationName: "UpdateTimeSpan",
		Variables: VariablesUpdate{
			OldStart: task.OldStart,
			Id:       task.Id,
			Start:    task.Start,
			Tags:     task.Tags,
			Note:     task.Note,
		},
		Query: "mutation UpdateTimeSpan($id: Int!, $start: Time!, $tags: [InputTimeSpanTag!], $note: String!) {\n  updateTimeSpan(id: $id, start: $start, tags: $tags, note: $note) {\n    id\n    start\n    tags {\n      key\n      value\n      __typename\n    }\n   note\n    __typename\n  }\n}\n",
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
	req.Header.Add("Cookie", t.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		fmt.Println("Updating task failure")
		fmt.Println(res.Status)
	}
}

func (t *Traggo) UpdateTimeSpanTask(task TimeSpanTask) {
	op := OperationUpdate{
		OperationName: "UpdateTimeSpan",
		Variables: VariablesUpdate{
			OldStart: task.OldStart,
			Id:       task.Id,
			Start:    task.Start,
			End:      task.End,
			Tags:     task.Tags,
			Note:     task.Note,
		},
		Query: "mutation UpdateTimeSpan($id: Int!, $start: Time!, $end: Time, $tags: [InputTimeSpanTag!], $oldStart: Time, $note: String!) {\n  updateTimeSpan(id: $id, start: $start, end: $end, tags: $tags, oldStart: $oldStart, note: $note) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
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
	req.Header.Add("Cookie", t.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		fmt.Println("Updating task failure")
		fmt.Println(res.Status)
	}
}

func (t *Traggo) Continue(task GenericTask) {
	op := OperationContinue{
		OperationName: "Continue",
		Variables: VariablesContinue{
			Id:    task.GetId(),
			Start: time.Now(),
		},
		Query: "mutation Continue($id: Int!, $start: Time!) {\n  copyTimeSpan(id: $id, start: $start) {\n    id\n    start\n    __typename\n  }\n}",
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
	req.Header.Add("Cookie", t.Token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != 200 {
		fmt.Println("Continuing task failure")
		fmt.Println(res.Status)
	}
}
