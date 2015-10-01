package cortana

import (
	"log"
	"bufio"
	"net"
	"sync"
	"bytes"
	"encoding/json"
)

var delim []byte = bytes.NewBufferString("\x0a").Bytes()

type SVConn struct {
	c net.Conn
	r *bufio.Reader 
	w *json.Encoder
}

func NewSVConn(c net.Conn) *SVConn {
	return &SVConn {
		c: c,
		r: bufio.NewReader(c),
		w: json.NewEncoder(c),
	}
}

func (s *SVConn) Read() (string, error) {
	return s.r.ReadString('\n')	
}

func (s *SVConn) Write(v interface {}) error {
	return s.w.Encode(v)
}

func (svc *SVConn) Close(s *Server) {
	s.remove_conn(svc)
	svc.c.Close()
}

func (svc *SVConn) RemoteAddr() net.Addr {
	return svc.c.RemoteAddr()
}

func (svc *SVConn) RemoteIP() net.IP {
	a := svc.RemoteAddr()
	ta, ok := a.(*net.TCPAddr)
	if !ok {
		log.Fatal("invalid address", a)
	}
	return ta.IP
}

type Response struct {
	Data Record
	Addr string
}
type Server struct {
	cmap map[string]*SVConn
	listener net.Listener
	ResponseCh chan Response
	mtx sync.Mutex
}

func NewServer(cnf Config) *Server {
	ln, err := net.Listen(cnf.BindAddr())
	if err != nil {
		log.Fatal(err)
	}
	s := Server {
		cmap: make(map[string]*SVConn), 
		listener: ln,
		ResponseCh: make(chan Response),
		mtx: sync.Mutex{},
	}
	return &s
}

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

func (s *Server) Send(payload interface {}) {
	for _,c := range s.cmap {
		c.Write(payload)
	}
}

func (s *Server) Close() {
	close(s.ResponseCh)
}

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
		s.ResponseCh <- Response{ Data: rec, Addr: svc.RemoteIP().String() }
	}
}

func (s *Server) add_conn(svc *SVConn) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.cmap[svc.RemoteAddr().String()] = svc
}

func (s *Server) remove_conn(svc *SVConn) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	delete(s.cmap, svc.RemoteAddr().String())
}
