package session

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (t *Traggo) ListBetweenDates(startDate time.Time, endDate time.Time) TimeSpanTaskList {
	op := OperationBetweenDate{
		OperationName: "TimeSpans",
		Variables: VariablesUpdateWithCursor{
			Start:  startDate,
			End:    endDate,
			Cursor: InputCursor{Offset: 0, PageSize: 100},
		},
		Query: "query TimeSpans($start: Time!, $end: Time!, $cursor: InputCursor) {\n  timeSpans(fromInclusive: $start, toInclusive: $end, cursor: $cursor) {\n    timeSpans {\n      id\n      start\n      end\n      tags {\n        key\n        value\n        __typename\n      }\n      oldStart\n      note\n      __typename\n    }\n    cursor {\n      hasMore\n    startId\n      offset\n      pageSize\n      __typename\n    }\n    __typename\n  }\n}\n",
	}
	timeSpanTaskSlice := []TimeSpanTask{}
	for {
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
			log.Fatal("List BetweenDates task failure")
		}

		var d TimeSpanRoot
		json.NewDecoder(res.Body).Decode(&d)
		for _, timespan := range d.Data.TimeSpans.TimeSpans {
			timeSpanTaskSlice = append(timeSpanTaskSlice, timespan)
		}
		if !d.Data.TimeSpans.Cursor.HasMore {
			// stop the pagination loop
			break
		}
		op.Variables.Cursor.Offset = d.Data.TimeSpans.Cursor.Offset

	}
	return timeSpanTaskSlice
}

func (t *Traggo) List() TimersData {
	op := OperationWithoutVariables{
		OperationName: "Trackers",
		Query:         "query Trackers {\n  timers {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
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
		log.Fatal("List task failure")
	}
	var d TimerTasks
	json.NewDecoder(res.Body).Decode(&d)

	return d.Data
}

// SearchTask by TaskID
// pagination do not support pageSize > 100
func (t *Traggo) SearchTask(id int) GenericTask {

	//Search for current running tasks (Trackers)
	for _, task := range t.List().Timers {
		if task.Id == id {
			return task
		}
	}
	//Search for old tasks
	op := OperationCursor{
		OperationName: "TimeSpans",
		Variables: VariablesCursor{
			Cursor: CursorRequest{Offset: 0, PageSize: 100},
		},
		Query: "query TimeSpans($cursor: InputCursor!) {\n  timeSpans(cursor: $cursor) {\n    timeSpans {\n      id\n      start\n      end\n      tags {\n        key\n        value\n        __typename\n      }\n      oldStart\n      note\n      __typename\n    }\n    cursor {hasMore\n      startId\n      offset\n      pageSize\n      __typename\n    }\n    __typename\n  }\n}\n",
	}

	for {
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
			log.Fatal("List task failure")
		}

		var d TimeSpanRoot
		json.NewDecoder(res.Body).Decode(&d)
		for _, task := range d.Data.TimeSpans.TimeSpans {
			if task.Id == id {
				return task
			}
		}
		if !d.Data.TimeSpans.Cursor.HasMore {
			// stop the pagination loop
			break
		}
		op.Variables.Cursor.Offset = d.Data.TimeSpans.Cursor.Offset
	}
	return nil
}

// SearchTask by TaskTag
// pagination do not support pageSize > 100
func (t *Traggo) SearchTaskByTag(tagName, tagValue string) GenericTask {

	//Search for current running tasks (Trackers)
	for _, task := range t.List().Timers {
		for _, taskTag := range task.Tags {
			if taskTag.Key == tagName && taskTag.Value == tagValue {
				return task
			}
		}
	}
	//Search for old tasks
	op := OperationCursor{
		OperationName: "TimeSpans",
		Variables: VariablesCursor{
			Cursor: CursorRequest{Offset: 0, PageSize: 100},
		},
		Query: "query TimeSpans($cursor: InputCursor!) {\n  timeSpans(cursor: $cursor) {\n    timeSpans {\n      id\n      start\n      end\n      tags {\n        key\n        value\n        __typename\n      }\n      oldStart\n      note\n      __typename\n    }\n    cursor {hasMore\n      startId\n      offset\n      pageSize\n      __typename\n    }\n    __typename\n  }\n}\n",
	}

	for {
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
			log.Fatal("List task failure")
		}

		var d TimeSpanRoot
		json.NewDecoder(res.Body).Decode(&d)
		for _, task := range d.Data.TimeSpans.TimeSpans {
			for _, taskTag := range task.Tags {
				if taskTag.Key == tagName && taskTag.Value == tagValue {
					return task
				}
			}
		}
		if d.Data.TimeSpans.Cursor.HasMore {
			op.Variables.Cursor.Offset = d.Data.TimeSpans.Cursor.Offset
		} else {
			// stop the pagination loop
			break
		}
	}
	return nil
}
