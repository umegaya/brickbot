package cortana_test

import (
    "testing"
	"../lib"
)

var line1 string = "{\"Kind\":\"hoge\", \"Pay}"
var line2 string = "{\"Kind\":\"hoge\", \"Payload\": {\"X\": 1, \"Y\": 2}}"

func TestRecordNew(t *testing.T) {
	r := cortana.NewRecord(line1)
	if r.Kind != "" || r.Payload != nil {
		t.Errorf("line parse should fails because of truncated line %v", r)
	}
	r2 := cortana.NewRecord(line2)
	if r2.Kind != "hoge" {
		t.Errorf("line parse should success %v", r2)
	}
	p, ok := r2.Payload.(map[string]interface {})
	if !ok || p["X"].(float64) != 1 || p["Y"].(float64) != 2 {
		t.Errorf("line parse should success %v", r2)		
	}
}
