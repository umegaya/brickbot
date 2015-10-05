package main

import (
	"log"
	"./lib"
)

func main() {
	var c cortana.Config
	err := c.Parse()
	if err != nil {
		log.Fatal(err)
	}
	sv := cortana.NewServer(c)
	go sv.Serv()
	cl := cortana.NewClient(c)
	cl.Run(sv, cortana.NewDockerController(c))
}
