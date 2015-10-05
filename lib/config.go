package cortana

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fsouza/go-dockerclient"
)

//ContainerConfig represents one container configuration.
//it is compatible for docker remote API (through go-dockerclient)
type ContainerConfig struct {
	Config     docker.Config     `json:"config"`
	HostConfig docker.HostConfig `json:"host_config"`
}

//DockerConfig represents entire docker config to create module containers
type DockerConfig struct {
	ServerAddress string                     `json:"server_address"`
	CertPath      string                     `json:"cert_path"`
	Containers    map[string]ContainerConfig `json:"module_containers"`
}

//Config represents entire configuration of slack-cortana
type Config struct {
	Token         string       `json:"token"`
	MainChannel   string       `json:"main_channel"`
	BindHost      string       `json:"bind_host"`
	BindPort      int          `json:"bind_port"`
	TemplatesPath string       `json:"templates_path"`
	Docker        DockerConfig `json:"docker"`
}

//ifip returns IP string of specified interface which name is *name*
func ifip(name string) net.IP {
	i, err := net.InterfaceByName(name)
	if err != nil {
		log.Fatal(err.Error())
	}
	l, err := i.Addrs()
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, a := range l {
		switch v := a.(type) {
		case *net.IPNet:
			return v.IP
		case *net.IPAddr:
			return v.IP
		}
	}
	return nil
}

//check_and_fill check configuration, if configuration seems not set, 
//it aborts or set default value
func (c *Config) check_and_fill() {
	if c.Token == "" {
		log.Fatal("config: token must be set")
	}
	if c.MainChannel == "" {
		c.MainChannel = "random"
	}
	if c.BindPort == 0 {
		c.BindPort = 8008
	}
	if c.TemplatesPath == "" {
		c.TemplatesPath = "./templates"
	}
	if c.Docker.CertPath == "" {
		log.Fatal("config: docker.cert_path must be set")
	}
	if c.Docker.ServerAddress == "localhost" {
		c.Docker.ServerAddress = ""
	}
	if c.BindHost == "" {
		var ip net.IP
		if c.Docker.ServerAddress == "" {
			ip = ifip("docker0")
		} else {
			ip = ifip("eth1")
		}
		if ip == nil {
			log.Fatal("network interface not properly configured")
		}
		c.BindHost = ip.String()
	}
	if c.Docker.ServerAddress == "" {
		ip := ifip("eth1")
		if ip == nil {
			log.Fatal("network interface not properly configured")
		}
		c.Docker.ServerAddress = ip.String()
	}
}

//Parse() pareses comannd line argument, and store it to newly created Config object, and return it.
func (c *Config) Parse() {
	s := flag.String("c", "", "configuration file for slack-cortana")
	flag.Parse()
	f, err := os.Open(*s)
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(f)
	if err := dec.Decode(c); err != nil {
		log.Fatal(err)
	}
	c.check_and_fill()
	log.Printf("network setting: %s:%d %s", c.BindHost, c.BindPort, c.Docker.ServerAddress)
}

//BindAddr returns bind address strings for net.Listen
func (c *Config) BindAddr() (string, string) {
	return "tcp", fmt.Sprintf("%s:%d", c.BindHost, c.BindPort)
}

