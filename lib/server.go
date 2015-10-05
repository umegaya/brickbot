package cortana

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"net"
	"sync"
)

var delim []byte = bytes.NewBufferString("\x0a").Bytes()

//SVConn repreeents one connection from module container
type SVConn struct {
	c net.Conn
	r *bufio.Reader
	w *json.Encoder
}

//NewSVConn creates new SVConn object and initialize by net.Conn (from server listener go routine)
func NewSVConn(c net.Conn) *SVConn {
	return &SVConn{
		c: c,
		r: bufio.NewReader(c),
		w: json.NewEncoder(c),
	}
}

//Read reads recieved message from SVConn
func (s *SVConn) Read() (string, error) {
	return s.r.ReadString('\n')
}

//Write writes json encodable object to SVConn
func (s *SVConn) Write(v interface{}) error {
	return s.w.Encode(v)
}

//Close closes SVConn. it removes connection from Server's address-connection map
func (svc *SVConn) Close(s *Server) {
	s.remove_conn(svc)
	svc.c.Close()
}

//RemoteAddr() returns net.Addr which associate with SVConn
func (svc *SVConn) RemoteAddr() net.Addr {
	return svc.c.RemoteAddr()
}

//RemoteIP() returns same as RemoteAddr but net.IP object
func (svc *SVConn) RemoteIP() net.IP {
	a := svc.RemoteAddr()
	ta, ok := a.(*net.TCPAddr)
	if !ok {
		log.Fatal("invalid address", a)
	}
	return ta.IP
}

//Response represents a response from connected module container
type Response struct {
	Data Record
	Addr string
}
//Server represents one server listener context.
type Server struct {
	cmap       map[string]*SVConn
	listener   net.Listener
	ResponseCh chan Response
	mtx        sync.Mutex
}

//NewServer() creates Server object from configuration
func NewServer(cnf Config) *Server {
	ln, err := net.Listen(cnf.BindAddr())
	if err != nil {
		log.Fatal(err)
	}
	s := Server{
		cmap:       make(map[string]*SVConn),
		listener:   ln,
		ResponseCh: make(chan Response),
		mtx:        sync.Mutex{},
	}
	return &s
}

//Serv() is main go routine of Server object. it accept connection from module container
//and run handler for each connection, as goroutine.
func (s *Server) Serv() {
Loop:
	for {
		c, err := s.listener.Accept()
		if err != nil {
			break Loop
		}
		log.Printf("accept from %s", c.RemoteAddr().String())
		go s.handler(NewSVConn(c))
	}
}

//Send sends json encodable payload object to each connected module containers
func (s *Server) Send(payload interface{}) {
	for _, c := range s.cmap {
		c.Write(payload)
	}
}

//Close stops Server object by closing channel and breaks Serv() goroutine.
func (s *Server) Close() {
	close(s.ResponseCh)
	s.listener.Close()
}

//handler reads record from module containers' connection, parse them, and send it to Client object's main goroutine
func (s *Server) handler(svc *SVConn) {
	s.add_conn(svc)
	defer svc.Close(s)
Loop:
	for {
		line, err := svc.Read()
		if err != nil {
			log.Printf("closed %s by %s", svc.RemoteAddr().String(), err.Error())
			break Loop
		}
		//log.Print("handler line = ", line)
		rec := NewRecord(line)
		//log.Printf("rec %v", rec)
		s.ResponseCh <- Response{Data: rec, Addr: svc.RemoteIP().String()}
	}
}

//add_conn adds Server's address-connection map to SVConn *svc*
func (s *Server) add_conn(svc *SVConn) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.cmap[svc.RemoteAddr().String()] = svc
}

//remove_conn removes SVConn *svc from Server's address-connection map.
func (s *Server) remove_conn(svc *SVConn) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	delete(s.cmap, svc.RemoteAddr().String())
}
