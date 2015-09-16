package cortana

import (
	"os"
	"log"
	"flag"
	"encoding/json"
)

type Config struct {
	Token string
}

func (c *Config) Parse() {
	s := flag.String("c", "", "configuration file for slack-cortana")
	flag.Parse()
	f, err := os.Open(*s); 
	if err != nil {
		log.Fatal(err)
	}
	dec := json.NewDecoder(f)
	if err := dec.Decode(c); err != nil {
		log.Fatal(err)
	}
}
