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

type ContainerConfig struct {
	Config     docker.Config     `json:"config"`
	HostConfig docker.HostConfig `json:"host_config"`
}

type DockerConfig struct {
	ServerAddress string                     `json:"server_address"`
	CertPath      string                     `json:"cert_path"`
	Containers    map[string]ContainerConfig `json:"module_containers"`
}

type Config struct {
	Token         string       `json:"token"`
	MainChannel   string       `json:"main_channel"`
	BindHost      string       `json:"bind_host"`
	BindPort      int          `json:"bind_port"`
	TemplatesPath string       `json:"templates_path"`
	Docker        DockerConfig `json:"docker"`
}

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

func (c *Config) BindAddr() (string, string) {
	return "tcp", fmt.Sprintf("%s:%d", c.BindHost, c.BindPort)
}
