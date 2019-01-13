package tlsrelay

import (
	"encoding/binary"
	"io"
	"net"
	"reflect"
	"sync"
	"time"
)

func (p *TLSRelay) handle(c net.Conn) {

	// Create connection to remote service
	dest, e := net.DialTCP("tcp", nil, p.rh)
	if e != nil {
		if p.isRelaying() {
			return
		}
	}

	period := time.Millisecond * 100

	e = c.(*net.TCPConn).SetKeepAlive(true)
	if e != nil {
		return
	}
	e = c.(*net.TCPConn).SetKeepAlivePeriod(period)
	if e != nil {
		return
	}

	if p.isRelaying() {
		e = dest.SetKeepAlive(true)
		if e != nil {
			return
		}
		e = dest.SetKeepAlivePeriod(period)
		if e != nil {
			return
		}
	}

	p.relay(c, dest)
}

func (p *TLSRelay) relay(src net.Conn, dst net.Conn) error {

	defer src.Close()
	if p.isRelaying() { // avoid crashing when it isn't initiated
		defer dst.Close()
	}

	// receive error from any go routing
	ret := make(chan error, 4)

	// go routines synchronization
	var wg sync.WaitGroup

	// prepare requests buffer (data from client to server)
	reqBuff := make(chan []byte, p.chBufferSize)

	// receive from client
	wg.Add(1)
	go p.connToChan(src, reqBuff, &wg, ret, src) // client = src

	if p.isRelaying() {
		// prepare response buffer (data from server to client)
		resBuff := make(chan []byte, p.chBufferSize)

		wg.Add(1)

		// send to server
		go p.chanToConn(reqBuff, dst, &wg, ret)

		// receive from server
		go p.connToChan(dst, resBuff, &wg, ret, src) // client = src

		// send to client
		go p.chanToConn(resBuff, src, &wg, ret)
	}

	// Check if connection is closed
	go p.checkIfClosed(src, &wg)
	go p.checkIfClosed(dst, &wg)

	wg.Wait()

	e := <-ret
	if e != io.EOF {
		return e
	}
	return nil
}

func (p *TLSRelay) chanToConn(src chan []byte, dst net.Conn, wg *sync.WaitGroup, res chan error) {
	defer wg.Done()

	var e error
	for {
		oneMsg := <-src
		msg := oneMsg
		for {
			if len(oneMsg) < p.packetSize {
				break // there can't be more bytes
			}

			time.Sleep(p.bufferTimeout)
			if len(src) == 0 {
				break // no more bytes after some ms
			}

			oneMsg := <-src
			msg = append(msg, oneMsg...)
		}

		_, e = dst.Write(msg)
		if e != nil {
			break
		}
		time.Sleep(p.sendDelay)
	}
	res <- e // should be nil or EOF
}

func (p *TLSRelay) connToChan(src net.Conn, dst chan []byte, wg *sync.WaitGroup, res chan error, client net.Conn) {
	defer wg.Done()

	var e error
loop:
	for {
		var n int

		src.SetReadDeadline(time.Now().Add(p.connTimeout)) // TODO: could we just exit when client or server close the connection?
		data := make([]byte, 66000)                        // Maximum TLS Record Size

		n, e = src.Read(data)
		if e != nil {
			break loop
		}

		for p.incompleteTLS(data[:n]) {
			newdata := make([]byte, 66000-len(data))
			m, e := src.Read(newdata)
			if e != nil {
				break loop
			}

			data = append(data[:n], newdata[:m]...)
			n += m
		}

		isRequest := reflect.DeepEqual(src, client)
		messages := p.manage(client.RemoteAddr().String(), isRequest, data[:n])
		for _, msg := range messages {
			dst <- msg
		}
		//dst <- data[:n]

		if !p.isRelaying() {
			break
		}
	}
	res <- e
}

func (p *TLSRelay) incompleteTLS(data []byte) bool {
	dataLen := len(data)
	tlsLen := int(binary.BigEndian.Uint16(data[3:5]))
	return (tlsLen > dataLen)
}

func (p *TLSRelay) checkIfClosed(c net.Conn, wg *sync.WaitGroup) {
	for {
		// Check if connections are closed
		break
	}
}
