package httpipe

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	wsChannelSize = 100
)

type wsData = string
type wsChannel = chan wsData

// Server type is an http.Server type with pre-defined responses
// for *.js pages, *.json pages, and other pages
type Server struct {
	srv         *http.Server
	mux         *http.ServeMux
	defaultPage string
	jsPage      string
	jsonPage    string

	upgrader websocket.Upgrader
	readWS   []wsChannel
	writeWS  []wsChannel
	mtxWS    sync.Mutex
}

// Handler function is handling HTTP requests and responding a
// different message depending on the extension (js, json or other)
func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	var res string

	url := r.RequestURI
	ext := filepath.Ext(url)
	switch ext {
	case ".js":
		res = s.jsPage
	case ".json":
		res = s.jsonPage
	default:
		res = s.defaultPage
	}

	resBytes := []byte(res)
	w.Write(resBytes)
}

// HandlerWS function is handling WebSocket connections
func (s *Server) handlerWS(w http.ResponseWriter, r *http.Request) {
	c, e := s.upgrader.Upgrade(w, r, nil)
	if e != nil {
		return
	}
	defer c.Close()

	i := s.createWS()
	go func() {
		ch := s.getWriteWS(i)
		for {
			data := <-ch
			msg := s.wsDataToByte(data)
			e := c.WriteMessage(websocket.TextMessage, msg)
			if e != nil {
				break
			}
		}
	}()

	ch := s.getReadWS(i)
	for {
		_, msg, e := c.ReadMessage()
		if e != nil {
			break
		}

		data := s.byteToWsData(msg)
		ch <- data
	}
}

// NewHTTPipe function creates a new HTTPipe webserver
func NewHTTPipe() Server {
	return *new(Server)
}

// StartBackground function starts a webserver from a goroutine
// addr describe where this webserver will listen for connections
func (s *Server) StartBackground(addr string) {
	s.mux = http.NewServeMux()
	s.srv = &http.Server{
		Addr:    addr,
		Handler: s.mux,
	}
	s.upgrader = websocket.Upgrader{}

	s.mux.HandleFunc("/", s.handler)
	s.mux.HandleFunc("/ws/", s.handlerWS)

	go func() {
		err := s.srv.ListenAndServe()
		if err == nil {
			log.Println("httpipe closes")
		}
	}()

	time.Sleep(time.Millisecond * 100)
}

// ChangeDefault function sets de default response
func (s *Server) ChangeDefault(p string) {
	s.defaultPage = p
}

// ChangeJS function sets the response for *.js requests
func (s *Server) ChangeJS(p string) {
	s.jsPage = p
}

// ChangeJSON function sets the response for *.json requests
func (s *Server) ChangeJSON(p string) {
	s.jsonPage = p
}

// HowManyWS function gets the amount of active WebSockers
func (s *Server) HowManyWS() int {
	return len(s.readWS)
}

func (s *Server) createWS() int {
	s.mtxWS.Lock()
	i := len(s.readWS)
	s.readWS = append(s.readWS, s.createWsChannel())
	s.writeWS = append(s.writeWS, s.createWsChannel())
	s.mtxWS.Unlock()
	return i
}

func (s *Server) getReadWS(i int) wsChannel {
	if i < len(s.readWS) {
		return s.readWS[i]
	}
	return nil
}

func (s *Server) getWriteWS(i int) wsChannel {
	if i < len(s.writeWS) {
		return s.writeWS[i]
	}
	return nil
}

func (s *Server) getWS(i int) (wsChannel, wsChannel) {
	return s.readWS[i], s.writeWS[i]
}

func (s *Server) createWsChannel() wsChannel {
	return make(wsChannel, wsChannelSize)
}

// ReadFromConn function reads from channel
func (s *Server) ReadFromConn(n int) wsData {
	ch := s.getReadWS(n)
	return <-ch
}

// WriteToConn function writes to channel
func (s *Server) WriteToConn(n int, data wsData) {
	ch := s.getWriteWS(n)
	ch <- data
}

func (s *Server) byteToWsData(msg []byte) wsData {
	return string(msg)
}

func (s *Server) wsDataToByte(data wsData) []byte {
	return []byte(data)
}

// StopBackground function shutdowns the webserver
func (s *Server) StopBackground() error {
	return s.srv.Shutdown(nil)
}
