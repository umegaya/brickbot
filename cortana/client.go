package cortana

import (
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

type Client struct {
	chmap map[string]string
	api *slack.Client
	rtm *slack.RTM	
}

func NewClient(cnf Config) (*Client, chan slack.RTMEvent) {
	var c Client
	c.api = slack.New(cnf.Token)
	//c.api.SetDebug(true)
	c.rtm = c.api.NewRTM()
	c.chmap = make(map[string]string)
	list, err := c.api.GetChannels(true)
	if err != nil {
		log.Fatal(err)
	}
	for _, elem := range list {
		fmt.Println("channels", elem.Name, elem.ID)
		c.chmap[elem.Name] = elem.ID
	}
	go c.rtm.ManageConnection()
	return &c, c.rtm.IncomingEvents
}

func(c *Client) Message(text, channel string) {
	cid, ok := c.chmap[channel]; 
	if !ok { 
		cid = channel 
	}
	c.rtm.SendMessage(c.rtm.NewOutgoingMessage(text, cid))
}
