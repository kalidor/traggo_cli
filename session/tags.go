package session

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

func (t tags) Contain(tagName string) bool {
	for _, tag := range t {
		if strings.EqualFold(tag.Key, tagName) {
			return true
		}
	}
	return false
}

func (t *Traggo) GetTags() tags {
	op := Operation{
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

func (t *Traggo) RemoveTag(tagName string) {
	variables := json.RawMessage(fmt.Sprintf(`{"key": "%s"}`, tagName))
	op := Operation{
		OperationName: "RemoveTag",
		Variables:     &variables,
		Query:         "mutation RemoveTag($key: String!) {\n  removeTag(key: $key) {\n    color\n    key\n  }\n}",
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
}
