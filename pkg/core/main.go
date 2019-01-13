package core

import (
	"sync"
	"time"

	"../httpipe"
	"../tlsrelay"
)

const (
	wildcard = "§1§"
)

// Core type is ...
type Core struct {
	// Servers
	pipe  httpipe.Server
	relay tlsrelay.TLSRelay

	// Heuristic Parameters
	minSize     int
	connN       int
	timeout     time.Duration
	timeoutConn time.Duration

	// Status parameters
	attack  uint8
	option  map[string]string
	running bool
	mtx     sync.Mutex
	charset []string

	// Output channels
	out    chan string
	error  chan error
	status chan string
	result chan string
}

// Exploit function...
func (c *Core) Exploit() (chan string, chan string, chan string, chan error) {
	// Set default parameters
	c.minSize = 0
	c.connN = 0
	c.timeout = time.Second * 4
	c.timeoutConn = time.Millisecond * 100

	// Prepare output
	c.mtx.Lock()
	if c.out == nil {
		c.out = make(chan string, 100)
	}
	if c.error == nil {
		c.error = make(chan error, 10)
	}
	if c.status == nil {
		c.status = make(chan string, 100)
	}
	if c.result == nil {
		c.result = make(chan string, 100)
	}
	c.running = true
	c.mtx.Unlock()

	go c.exploit()
	return c.status, c.out, c.result, c.error
}

// Break function stop the execution of Exploit function
func (c *Core) Break() {
	// exit if already stoped
	if !c.running {
		return
	}

	// set state to not running
	c.mtx.Lock()
	c.running = false
	c.mtx.Unlock()

	// clean status and error channels
	for len(c.status) > 0 {
		<-c.status
	}
	for len(c.error) > 0 {
		<-c.error
	}

	// stop services
	c.status <- "Shutdown Control HTTP server"
	c.pipe.StopBackground()
	c.status <- "Shutdown TLS Relay server"
	c.relay.Stop()

	// close output channel, force to continue
	c.mtx.Lock()
	if c.out != nil {
		close(c.out)
		c.out = nil
	}
	if c.status != nil {
		close(c.status)
		c.status = nil
	}
	if c.error != nil {
		close(c.error)
		c.error = nil
	}
	c.mtx.Unlock()
}
