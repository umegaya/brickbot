package main

import (
	"./lib"
)

func main() {
	var c cortana.Config
	c.Parse()
	sv := cortana.NewServer(c)
	go sv.Serv()
	cl := cortana.NewClient(c)
	cl.Run(sv)
}

