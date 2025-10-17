package session

import (
	"log"
)

type tag struct {
	Key    string `json:"key"`
	Usages int    `json:"usages"`
}
type tags []tag

type datatags struct {
	Tags tags `json:"tags"`
}
type rootTags struct {
	Data datatags `json:"data"`
}

func (t *Traggo) GetTags() tags {
	op := OperationWithoutVariables{
		OperationName: "Tags",
		Query:         "query Tags {\n  tags {\n    key\n    usages\n}\n}",
	}

	// Parse http.Response Boby as JSON and display it
	var d rootTags
	err := t.Request(
		"GetTags",
		"POST",
		op,
		&d,
	)
	if err != nil {
		log.Fatal(err)
	}
	return d.Data.Tags

}
