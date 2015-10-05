package cortana_test

import (
    "testing"
	"../lib"

	"github.com/fsouza/go-dockerclient"
)

var dummy_conf2 cortana.Config = cortana.Config{
	Token: "hoge",
	TemplatesPath: "./templates",
	BindPort: 8888,
	BindHost: "0.0.0.0",
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

func TestServerNew(t *testing.T) {
	s := cortana.NewServer(dummy_conf2)
	s.Close()
}
