package cortana_test

import (
    "testing"
    "strings"
	"../lib"
)

func TestConfigLoad(t *testing.T) {
	var c cortana.Config
	if err := c.Load("./fixture/settings/not_exist.json"); err == nil || !strings.Contains(err.Error(), "no such file") {
		t.Errorf("load should fail because of lack of file")
	}
	if err := c.Load("./fixture/settings/broken.json"); err == nil || !strings.Contains(err.Error(), "EOF") {
		t.Errorf("load should fail because of broken json file %v", err)
	}	
	if err := c.Load("./fixture/settings/no_token.json"); err == nil || !strings.Contains(err.Error(), "token"){
		t.Errorf("load should fail because of no token %v", err)
	}
	if err := c.Load("./fixture/settings/no_cert_path.json"); err == nil || !strings.Contains(err.Error(), "docker.cert_path") {
		t.Errorf("load should fail because of no cert path %v", err)
	}
	if err := c.Load("./fixture/settings/try_fill.json"); err != nil {
		t.Errorf("load should success but %v", err)
	}
	if c.BindPort != 8008 {
		t.Errorf("BindPort: default value wrong %v", c.BindPort)
	}
	if c.TemplatesPath != "./templates" {
		t.Errorf("BindPort: default value wrong %v", c.TemplatesPath)
	}
	if c.MainChannel != "random" {
		t.Errorf("BindPort: default value wrong %v", c.MainChannel)
	}
}
