package session

import (
	"log"
	"strings"
	"time"

	"github.com/kalidor/traggo_cli/config"
)

var TimeNow = time.Now

type createTimeSpanData struct {
	Data TimerTask `json:"createTimeSpan"`
}
type createTimeSpanRoot struct {
	Data   createTimeSpanData `json:"data"`
	Errors []Error            `json:"errors"`
}

func (t *Traggo) Start(tags []string, note string) {
	var genTags []Tag
	for _, tag := range tags {
		if strings.Contains(tag, ":") {
			s := strings.SplitN(tag, ":", 2)
			genTags = append(genTags, Tag{Key: s[0], Value: s[1]})
		}
	}

	variables := struct {
		Start time.Time `json:"start"`
		Tags  []Tag     `json:"tags"`
		Note  string    `json:"note"`
	}{
		Start: TimeNow().UTC(),
		Tags:  genTags,
		Note:  note,
	}

	op := Operation{
		OperationName: "StartTimer",
		Variables:     variables,
		Query:         "mutation StartTimer($start: Time!, $tags: [InputTimeSpanTag!], $note: String!) {\n  createTimeSpan(start: $start, tags: $tags, note: $note) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
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
	d.Data.Data.PreparePretty(t.Colors)

}

func (t *Traggo) Stop(colors config.ColorsDef, ids []int) {
	variables := struct {
		Id  int       `json:"id"`
		End time.Time `json:"end"`
	}{
		Id:  0,
		End: TimeNow().UTC().Add(time.Hour * 1),
	}
	op := Operation{
		OperationName: "StopTimer",
		Variables:     variables,
		Query:         "mutation StopTimer($id: Int!, $end: Time!) {\n  stopTimeSpan(id: $id, end: $end) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
	}

	for _, id := range ids {
		variables.Id = id
		variables.End = TimeNow().UTC().Add(time.Hour * 1)

		op.Variables = &variables
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
		d.PreparePretty(colors)
	}
}

func (t *Traggo) Delete(ids []int) {
	variables := struct {
		Id int `json:"id"`
	}{
		Id: 0,
	}
	op := Operation{
		OperationName: "RemoveTimeSpan",
		Variables:     variables,
		Query:         "mutation RemoveTimeSpan($id: Int!) {\n  removeTimeSpan(id: $id) {\n    id\n    __typename\n  }\n}\n",
	}
	for _, id := range ids {
		variables.Id = id
		op.Variables = variables
		t.Request("RemoveTimeSpan", "POST", op, nil)
	}
}

func (t *Traggo) UpdateTimerTask(task TimerTask) {
	variables := struct {
		OldStart time.Time `json:"oldStart,omitzero"`
		Id       int       `json:"id,omitempty"`
		Start    time.Time `json:"start,omitzero"`
		Tags     []Tag     `json:"tags,omitzero"`
		Note     string    `json:"note"` // do not omit if empty
	}{
		OldStart: task.OldStart,
		Id:       task.Id,
		Start:    task.Start,
		Tags:     task.Tags,
		Note:     task.Note,
	}

	op := Operation{
		OperationName: "UpdateTimeSpan",
		Variables:     variables,
		Query:         "mutation UpdateTimeSpan($id: Int!, $start: Time!, $tags: [InputTimeSpanTag!], $note: String!) {\n  updateTimeSpan(id: $id, start: $start, tags: $tags, note: $note) {\n    id\n    start\n    tags {\n      key\n      value\n      __typename\n    }\n   note\n    __typename\n  }\n}\n",
	}
	t.Request("UpdateTimeSpan", "POST", op, nil)

}

func (t *Traggo) UpdateTimeSpanTask(task TimeSpanTask) {
	variables := struct {
		OldStart time.Time `json:"oldStart,omitzero"`
		Id       int       `json:"id,omitempty"`
		Start    time.Time `json:"start,omitzero"`
		End      time.Time `json:"end,omitzero"`
		Tags     []Tag     `json:"tags,omitzero"`
		Note     string    `json:"note"` // do not omit if empty
	}{
		OldStart: task.OldStart,
		Id:       task.Id,
		Start:    task.Start,
		End:      task.End,
		Tags:     task.Tags,
		Note:     task.Note,
	}
	op := Operation{
		OperationName: "UpdateTimeSpan",
		Variables:     variables,
		Query:         "mutation UpdateTimeSpan($id: Int!, $start: Time!, $end: Time, $tags: [InputTimeSpanTag!], $oldStart: Time, $note: String!) {\n  updateTimeSpan(id: $id, start: $start, end: $end, tags: $tags, oldStart: $oldStart, note: $note) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
	}
	t.Request("UpdateTimeSpanTask", "POST", op, nil)
}

func (t *Traggo) Continue(task GenericTask) {
	variables := struct {
		Id    int       `json:"id,omitempty"`
		Start time.Time `json:"start"`
	}{
		Id:    task.GetId(),
		Start: TimeNow(),
	}
	op := Operation{
		OperationName: "Continue",
		Variables:     variables,
		Query:         "mutation Continue($id: Int!, $start: Time!) {\n  copyTimeSpan(id: $id, start: $start) {\n    id\n    start\n    __typename\n  }\n}",
	}
	t.Request("Continue", "POST", op, nil)

}
