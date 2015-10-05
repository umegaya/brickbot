package cortana_test

import (
	"encoding/json"
    "testing"
    "log"
	"../lib"

	"github.com/fsouza/go-dockerclient"
)

var dummy_conf cortana.Config = cortana.Config{
	Token: "hoge",
	TemplatesPath: "./templates",
	Docker: cortana.DockerConfig {
		Containers: map[string]cortana.ContainerConfig {
			"fuga": cortana.ContainerConfig{
				Config: docker.Config {
					Image: "foo/bar",
				},
			},
		},
	},
}

type test_sub_payload struct {
	A string
}
type test_payload struct {
	A string
	X test_sub_payload
}

func to_s(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(b)
}

func NewDummyClient() *cortana.Client {
	return cortana.NewClient(dummy_conf)
}

//TestClientTemplate tests client can load template correctly and apply to input payload
func TestClientTemplate(t *testing.T) {
	c := NewDummyClient()
	c.LoadTemplate("./fixture/templates", "guha")
	if 0 != len(c.Templates()) {
		t.Errorf("template should not be loaded because of wrong name %d", len(c.Templates()))
	}
	c.LoadTemplate("./fixture/templates", "fuga")
	if 0 == len(c.Templates()) {
		t.Errorf("template should be loaded %d", len(c.Templates()))		
	}
	p := test_payload {
		A: "a",
		X: test_sub_payload {
			A: "xa",
		},
	}
	txt := c.FormatMessage("guha", "Example", p)
	if txt != to_s(p) {
		t.Errorf("message not formatted because of no such container %s, %s", txt, to_s(p))
	}
	txt = c.FormatMessage("fuga", "NotExample", p)
	if txt != to_s(p) {
		t.Errorf("message not formatted because of no such message kind %s, %s", txt, to_s(p))
	}
	txt = c.FormatMessage("fuga", "Example", p)
	if txt != "hi, a = a, x.a=xa" {
		t.Errorf("message should be formatted %s, %s", txt, to_s(p))
	}
}
