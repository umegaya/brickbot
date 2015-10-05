package cortana

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/nlopes/slack"
)

//Client represents running context of one slack-cortana instance
type Client struct {
	chmap     map[string]string
	templates map[string]map[string]*template.Template
	api       *slack.Client
	cnf       Config
	sig       chan os.Signal
	closer    chan interface{}
	dc        *DockerController
	rtm       *slack.RTM
}

//NewClient create and initialize Client object by Config *cnf*
func NewClient(cnf Config) *Client {
	var c Client
	c.cnf = cnf
	c.api = slack.New(cnf.Token)
	//c.api.SetDebug(true)
	c.rtm = c.api.NewRTM()
	c.sig = make(chan os.Signal)
	c.closer = make(chan interface{})
	c.chmap = make(map[string]string)
	c.templates = make(map[string]map[string]*template.Template)
	return &c
}

//Initialize initializes Client object 
func (c *Client) Initialize(dc *DockerController) *Client {
	c.dc = dc
	list, err := c.api.GetChannels(true)
	if err != nil {
		log.Fatal(err)
	}
	for _, elem := range list {
		fmt.Println("channels", elem.Name, elem.ID)
		c.chmap[elem.Name] = elem.ID
	}
	for name, _ := range c.dc.Containers {
		c.LoadTemplate(c.cnf.TemplatesPath, name)
	}
	go c.rtm.ManageConnection()
	return c
}

//LoadTemplate load and store text/template object into internal map object Client.templates.
//path and name are given from configuration. note that name is a key of module-containers configuration.
//so you have to put *name*.json in *path* directory to load template for container which configuration name is *name*.
func (c *Client) LoadTemplate(path, name string) {
	fullpath := fmt.Sprintf("%s/%s.json", path, name)
	log.Printf("fullpath:%s", fullpath)
	f, err := os.Open(fullpath)
	if err != nil {
		return //ok. no template created
	}
	var tmpm map[string]string
	dec := json.NewDecoder(f)
	if err := dec.Decode(&tmpm); err != nil {
		log.Fatal(err)
	}
	r := make(map[string]*template.Template)
	for n, body := range tmpm {
		tpl := template.New(fmt.Sprintf("%s.%s", name, n))
		tpl.Parse(body)
		r[n] = tpl
	}
	c.templates[name] = r
}

//Templates returns templates map of Client object. only for test use
func (c *Client) Templates() map[string]map[string]*template.Template {
	return c.templates
}

//close_watcher wait container stops, and send signal to main go routine (Client.Run)
func (c *Client) close_watcher() {
	signal.Notify(c.sig,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go (func() {
		sig := <-c.sig
		c.closer <- sig
	})()
}

//Run is main go routine of Client object it receives RTMEvent from slack object and send it to connected client,
//also receives message from connected client, and send it to config.MainChannel
func (c *Client) Run(sv *Server, dc *DockerController) {
	c.Initialize(dc)
	defer sv.Close()
	defer c.dc.Stop()
	c.close_watcher()
Loop:
	for {
		select {
		case msg := <-c.rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				fmt.Println("Infos:", ev.Info)
				//c.Message("I'm up now", c.cnf.MainChannel)

			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)

			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// Ignore other events..
				fmt.Printf("Unexpected: %v\n", msg.Data)
			}
			sv.Send(msg)
		case resp := <-sv.ResponseCh:
			ct, name := c.dc.FindContainer(resp.Addr)
			if ct != nil {
				txt := c.FormatMessage(name, resp.Data.Kind, resp.Data.Payload)
				c.Message(txt, c.cnf.MainChannel)
			}
		case sig := <-c.closer:
			log.Printf("singal recieved: %d", sig)
			break Loop
		}
	}
}

//FormatMessage formats human friendly message from payload which received from connected module-containers.
func (c *Client) FormatMessage(msg_from, msg_kind string, payload interface{}) string {
	entries, ok := c.templates[msg_from]
	if !ok {
		b, err := json.Marshal(payload)
		if err != nil {
			return err.Error()
		}
		return string(b)
	}
	tpl, ok := entries[msg_kind]
	if !ok {
		b, err := json.Marshal(payload)
		if err != nil {
			return err.Error()
		}
		return string(b)
	}
	b := new(bytes.Buffer)
	tpl.Execute(b, payload)
	return b.String()
}

//Message implements SlackClient interface
func (c *Client) Message(text, channel string) {
	cid, ok := c.chmap[channel]
	if !ok {
		cid = channel
	}
	c.rtm.SendMessage(c.rtm.NewOutgoingMessage(text, cid))
}
