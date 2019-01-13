package tlsrelay

import (
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/google/gopacket"
)

// TLSrelay custom errors
const (
	noRelayError      = "Relay is not running"
	relayRunningError = "Relay is still running"

	portError = "Invalid port"
	hostError = "Invalid hostname or IP address"

	defaultHost = "127.0.0.1"
	defaultPort = 8080

	defaultPacketSize        = 1500
	defaultChBufferSize      = 10
	defaultBufferTimeout     = time.Millisecond * 1
	defaultConnectionTimeout = time.Millisecond * 5000
	defaultWaitToStartTime   = time.Millisecond * 100
	defaultSendDelay         = 0
	maxSendDelay             = time.Millisecond * 100

	defaultTLSInfoChBufferSize = 100
)

// TLSRelay structure...
type TLSRelay struct {
	listenHost string
	listenPort int
	relayHost  string
	relayPort  int
	lh         *net.TCPAddr
	rh         *net.TCPAddr
	L          *net.TCPListener
	stop       bool
	closeError error

	packetSize      int
	chBufferSize    int
	bufferTimeout   time.Duration
	connTimeout     time.Duration
	waitToStartTime time.Duration
	sendDelay       time.Duration

	tlsInfo       chan string
	tlsInfoChSize int

	gopacketOpt gopacket.DecodeOptions
}

// NewTLSRelay function...
func NewTLSRelay() TLSRelay {
	p := TLSRelay{}

	p.listenHost = defaultHost
	p.listenPort = defaultPort
	p.stop = true
	p.closeError = errors.New(noRelayError)

	p.packetSize = defaultPacketSize
	p.chBufferSize = defaultChBufferSize
	p.bufferTimeout = defaultBufferTimeout
	p.connTimeout = defaultConnectionTimeout
	p.waitToStartTime = defaultWaitToStartTime
	p.sendDelay = defaultSendDelay

	p.tlsInfoChSize = defaultTLSInfoChBufferSize
	p.tlsInfo = make(chan string, p.tlsInfoChSize)

	p.gopacketOpt = gopacket.DecodeOptions{
		SkipDecodeRecovery:       true,
		DecodeStreamsAsDatagrams: true,
	}

	return p
}

// Start function ...
func (p *TLSRelay) Start() error {
	addr := p.listenHost + ":" + strconv.Itoa(p.listenPort)
	lh, e1 := net.ResolveTCPAddr("tcp", addr)
	if e1 != nil {
		return e1
	}
	p.lh = lh

	addr = p.relayHost + ":" + strconv.Itoa(p.relayPort)
	if addr != ":0" {
		rh, e2 := net.ResolveTCPAddr("tcp", addr)
		if e2 != nil {
			return e2
		}
		p.rh = rh
	}

	l, e3 := net.ListenTCP("tcp", lh)
	if e3 != nil {
		return e3
	}
	p.L = l
	p.closeError = errors.New(relayRunningError)
	p.stop = false

	for !p.stop {
		conn, e4 := p.L.Accept()
		if e4 != nil {
			continue
		}
		go p.handle(conn)
	}

	return nil
}

// WaitToStart function...
func (p *TLSRelay) WaitToStart() {
	time.Sleep(p.waitToStartTime)
}

// Stop function...
func (p *TLSRelay) Stop() error {
	// Return error if Proxy was not initialized
	if p.isNotInitialized() || (p.L == nil) {
		return errors.New(noRelayError)
	}
	// This should stop the loop in Start()
	p.stop = true
	// Close the listener
	p.closeError = p.L.Close()
	return p.closeError
}

// SetListen function...
func (p *TLSRelay) SetListen(host string, port int) error {
	e := p.checkHostPort(host, port)
	if e != nil {
		return e
	}

	p.listenHost = host
	p.listenPort = port
	return nil
}

// SetRelay function...
func (p *TLSRelay) SetRelay(host string, port int) error {
	e := p.checkHostPort(host, port)
	if e != nil {
		return e
	}

	p.relayHost = host
	p.relayPort = port
	return nil
}

func (p *TLSRelay) checkHostPort(host string, port int) error {
	if (port <= 0) || (port > 65535) {
		return errors.New(portError)
	}

	if net.ParseIP(host) != nil {
		return nil
	}

	_, e := net.ResolveIPAddr("ip", host)
	return e
}

func (p *TLSRelay) isNotInitialized() bool {
	l := len(p.listenHost) // noinit -> ""
	v := p.listenPort      // noinit -> 0
	return (l+v == 0)
}

func (p *TLSRelay) isRelaying() bool {
	return !p.isNotRelaying()
}

func (p *TLSRelay) isNotRelaying() bool {
	l := len(p.relayHost) // noinit -> ""
	v := p.relayPort      // noinit -> 0
	return (l+v == 0)
}

// ResetRelay function...
func (p *TLSRelay) ResetRelay() {
	p.relayHost = defaultHost
	p.relayPort = defaultPort
}

// Demultiplex function...
func (p *TLSRelay) Demultiplex(x bool) {
	if x {
		p.sendDelay = maxSendDelay
	} else {
		p.sendDelay = defaultSendDelay
	}
}
