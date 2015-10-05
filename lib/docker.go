package cortana

import (
	"fmt"
	"log"

	"github.com/fsouza/go-dockerclient"
)

//DockerController represents context of running container, started by slack-cortana.
type DockerController struct {
	Containers map[string]*docker.Container
	client     *docker.Client
}

//NewDockerController creates DockerController object and initializing them by running container with given configuration.
func NewDockerController(c Config) *DockerController {
	l := len(c.Docker.Containers)
	if l <= 0 {
		return &DockerController{
			Containers: make(map[string]*docker.Container),
		}
	}
	ret := make(map[string]*docker.Container, l)
	path := c.Docker.CertPath
	ca := fmt.Sprintf("%s/ca.pem", path)
	cert := fmt.Sprintf("%s/cert.pem", path)
	key := fmt.Sprintf("%s/key.pem", path)
	client, err := docker.NewTLSClient(fmt.Sprintf("tcp://%s:2376", c.Docker.ServerAddress), cert, key, ca)
	if err != nil {
		log.Fatal(err)
	}
	_, addr := c.BindAddr()
	for name, cnf := range c.Docker.Containers {
		rmopts := docker.RemoveContainerOptions{
			ID:            name,
			RemoveVolumes: true,
			Force:         true,
		}
		if err := client.RemoveContainer(rmopts); err != nil {
			log.Print("remove container fails:", err)
		}
		cnf.Config.Env = append(cnf.Config.Env, fmt.Sprintf("BRICKBOT_ADDR=%s", addr))
		opts := docker.CreateContainerOptions{
			Name:       name,
			Config:     &cnf.Config,
			HostConfig: &cnf.HostConfig,
		}
		ct, err := client.CreateContainer(opts)
		if err != nil {
			log.Fatal("create container fails:", err)
		}
		if err := client.StartContainer(ct.ID, &cnf.HostConfig); err != nil {
			log.Fatal("start container fails:", err)
		}
		ret[name] = ct
	}
	return &DockerController{
		Containers: ret,
		client:     client,
	}
}

//Stop stops all container started by this program, its called via defer in client.Run()
func (dc *DockerController) Stop() {
	for _, c := range dc.Containers {
		log.Print("kill ct", c.ID)
		dc.client.KillContainer(docker.KillContainerOptions{
			ID: c.ID,
		})
	}
}

//FindContainer finds container from tcp connection address. returns nil if no container found.
func (dc *DockerController) FindContainer(addr string) (*docker.Container, string) {
	for name, c := range dc.Containers {
		if c.NetworkSettings == nil {
			tmp, err := dc.client.InspectContainer(c.ID)
			if err != nil {
				log.Fatal("inspect container fails:", err)
			}
			c.NetworkSettings = tmp.NetworkSettings
		}
		if c.NetworkSettings.IPAddress == addr {
			return c, name
		}
	}
	return nil, ""
}
