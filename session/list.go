package session

import (
	"log"
	"time"
)

// ListBetweenDates return []TimeSpanTask{} containing matching tasks
// between the provided dates
func (t *Traggo) ListBetweenDates(startDate time.Time, endDate time.Time) TimeSpanTaskList {
	op := OperationBetweenDate{
		OperationName: "TimeSpans",
		Variables: VariablesUpdateWithCursor{
			Start:  startDate,
			End:    endDate,
			Cursor: CursorRequest{Offset: 0, PageSize: 10},
		},
		Query: "query TimeSpans($start: Time!, $end: Time!, $cursor: InputCursor) {\n  timeSpans(fromInclusive: $start, toInclusive: $end, cursor: $cursor) {\n    timeSpans {\n      id\n      start\n      end\n      tags {\n        key\n        value\n        __typename\n      }\n      oldStart\n      note\n      __typename\n    }\n    cursor {\n      hasMore\n    startId\n      offset\n      pageSize\n      __typename\n    }\n    __typename\n  }\n}\n",
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
		op.Variables.Cursor.Offset = d.Data.TimeSpans.Cursor.Offset

	}
	return timeSpanTaskSlice
}

// ListCurrentTasks return TimerTasks containing current running tasks
func (t *Traggo) ListCurrentTasks() TimersData {
	op := OperationWithoutVariables{
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
	op := OperationCursor{
		OperationName: "TimeSpans",
		Variables: VariablesCursor{
			Cursor: CursorRequest{Offset: 0, PageSize: 100},
		},
		Query: "query TimeSpans($cursor: InputCursor!) {\n  timeSpans(cursor: $cursor) {\n    timeSpans {\n      id\n      start\n      end\n      tags {\n        key\n        value\n        __typename\n      }\n      oldStart\n      note\n      __typename\n    }\n    cursor {hasMore\n      startId\n      offset\n      pageSize\n      __typename\n    }\n    __typename\n  }\n}\n",
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
		op.Variables.Cursor.Offset = d.Data.TimeSpans.Cursor.Offset
	}
	return tasks
}
