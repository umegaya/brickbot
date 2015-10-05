package cortana_test

import (
    "testing"
	"../lib"

	"github.com/fsouza/go-dockerclient"
)

var dc cortana.DockerController = cortana.DockerController {
	Containers: map[string]*docker.Container {
		"c1": {
			NetworkSettings: &docker.NetworkSettings {
				IPAddress: "1.1.1.1",
			},
		},
		"c2": {
			NetworkSettings: &docker.NetworkSettings {
				IPAddress: "2.2.2.2",
			},
		},
	},
}

func TestDockerControllerFindContainer(t *testing.T) {
	var c *docker.Container
	var name string
	c, name = dc.FindContainer("3.3.3.3")
	if c != nil {
		t.Errorf("container should not find %v", name)
	}
	c, name = dc.FindContainer("2.2.2.2")
	if c == nil || c.NetworkSettings.IPAddress != "2.2.2.2" {
		t.Errorf("container should find and correct ipaddr %v", c)
	}
}
