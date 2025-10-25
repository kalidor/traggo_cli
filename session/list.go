package session

import (
	"log"
	"time"
)

// ListBetweenDates return []TimeSpanTask{} containing matching tasks
// between the provided dates
func (t *Traggo) ListBetweenDates(startDate time.Time, endDate time.Time) TimeSpanTaskList {
	variables := struct {
		Start  time.Time     `json:"start"`
		End    time.Time     `json:"end"`
		Cursor CursorRequest `json:"cursor"`
	}{
		Start:  startDate,
		End:    endDate,
		Cursor: CursorRequest{Offset: 0, PageSize: 10},
	}

	op := Operation{
		OperationName: "TimeSpans",
		Variables:     variables,
		Query:         "query TimeSpans($start: Time!, $end: Time!, $cursor: InputCursor) {\n  timeSpans(fromInclusive: $start, toInclusive: $end, cursor: $cursor) {\n    timeSpans {\n      id\n      start\n      end\n      tags {\n        key\n        value\n        __typename\n      }\n      oldStart\n      note\n      __typename\n    }\n    cursor {\n      hasMore\n    startId\n      offset\n      pageSize\n      __typename\n    }\n    __typename\n  }\n}\n",
	}
	timeSpanTaskSlice := []TimeSpanTask{}
	for {

		var d TimeSpanRoot
		err := t.Request(
			"ListBetweenDates",
			"POST",
			op,
			&d,
		)
		if err != nil {
			log.Fatal(err)
		}

		for _, timespan := range d.Data.TimeSpans.TimeSpans {
			timeSpanTaskSlice = append(timeSpanTaskSlice, timespan)
		}
		if !d.Data.TimeSpans.Cursor.HasMore {
			// stop the pagination loop
			break
		}
		op.Variables = struct {
			Start  time.Time     `json:"start"`
			End    time.Time     `json:"end"`
			Cursor CursorRequest `json:"cursor"`
		}{
			Start:  startDate,
			End:    endDate,
			Cursor: CursorRequest{Offset: d.Data.TimeSpans.Cursor.Offset, PageSize: 10},
		}

	}
	return timeSpanTaskSlice
}

// ListCurrentTasks return TimerTasks containing current running tasks
func (t *Traggo) ListCurrentTasks() TimersData {
	op := Operation{
		OperationName: "Trackers",
		Query:         "query Trackers {\n  timers {\n    id\n    start\n    end\n    tags {\n      key\n      value\n      __typename\n    }\n    oldStart\n    note\n    __typename\n  }\n}\n",
	}
	var tasks TimerTasks
	err := t.Request(
		"ListCurrentTasks",
		"POST",
		op,
		&tasks,
	)
	if err != nil {
		log.Fatal(err)
	}

	return tasks.Data
}

func (t *Traggo) ListCompleteTasks() TimeSpanTaskList {
	variables := struct {
		Cursor CursorRequest `json:"cursor"`
	}{
		Cursor: CursorRequest{Offset: 0, PageSize: 100},
	}
	op := Operation{
		OperationName: "TimeSpans",
		Variables:     variables,
		Query:         "query TimeSpans($cursor: InputCursor!) {\n  timeSpans(cursor: $cursor) {\n    timeSpans {\n      id\n      start\n      end\n      tags {\n        key\n        value\n        __typename\n      }\n      oldStart\n      note\n      __typename\n    }\n    cursor {hasMore\n      startId\n      offset\n      pageSize\n      __typename\n    }\n    __typename\n  }\n}\n",
	}

	var tasks TimeSpanTaskList

	for {
		var d TimeSpanRoot
		err := t.Request(
			"ListCompleteTasks",
			"POST",
			op,
			&d,
		)
		if err != nil {
			log.Fatal(err)
		}

		tasks = append(tasks, d.Data.TimeSpans.TimeSpans...)

		// stop the pagination loop
		if !d.Data.TimeSpans.Cursor.HasMore {
			break
		}
		op.Variables = struct {
			Cursor CursorRequest `json:"cursor"`
		}{
			Cursor: CursorRequest{Offset: d.Data.TimeSpans.Cursor.Offset, PageSize: 100},
		}
	}
	return tasks
}
