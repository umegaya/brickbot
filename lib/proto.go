package cortana

import (
	"encoding/json"
	"log"
	"strings"
)

/*
	[
		"Kind": "message_kind",
		"Payload": {
			...
		}
	]
*/

type Record struct {
	Kind    string      `json:"Kind"`
	Payload interface{} `json:"Payload"`
}

func NewRecord(line string) Record {
	var r Record
	err := json.NewDecoder(strings.NewReader(line)).Decode(&r)
	if err != nil {
		log.Printf("newrec err", err.Error())
		return Record{}
	}
	return r
}
