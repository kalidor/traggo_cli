package session

import (
	"log"
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
			Start: time.Now().UTC().Add(time.Hour * 2),
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
	d.Data.Data.PreparePretty(t.Colors)

}

func (t *Traggo) Stop(colors config.ColorsDef, ids []int) {
	op := OperationUpdate{
		OperationName: "StopTimer",
		Variables: VariablesUpdate{
			Id:  0,
			End: time.Now().UTC().Add(time.Hour * 2),
		},
		Query: "mutation StopTimer($id: Int!, $end: Time!) {\n  stopTimeSpan(id: $id, end: $end) {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
	}

	for _, id := range ids {
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
		d.PreparePretty(colors)
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
		t.Request("RemoveTimeSpan", "POST", op, nil)
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
	t.Request("UpdateTimeSpan", "POST", op, nil)

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
	t.Request("UpdateTimeSpanTask", "POST", op, nil)
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
	t.Request("Continue", "POST", op, nil)

}
