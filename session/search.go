package session

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// SearchTask by TaskID in current running tasks and already done tasks.
func (t *Traggo) SearchTask(id int) GenericTask {

	//Search for current running tasks (Trackers)
	for _, task := range t.ListCurrentTasks().Timers {
		if task.Id == id {
			return task
		}
	}
	//Search for old tasks
	all := t.ListCompleteTasks()
	for _, task := range all {
		if task.Id == id {
			return task
		}
	}
	return nil
}

// SearchTaskByTag look for task matching provided tagName and tagValue.
func (t *Traggo) SearchTaskByTag(tagName, tagValue string) GenericTask {

	//Search for current running tasks (Trackers)
	for _, task := range t.ListCurrentTasks().Timers {
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
