package tlsrelay

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestTLSRelayRelay(t *testing.T) {
	// Prepare tests
	tests := []struct {
		name       string
		content    []byte
		isRequest  bool
		isResponse bool
	}{
		{"Request four bytes", []byte{0xBE, 0xEF, 0xCA, 0xFE}, true, false},
		{"Response four bytes", []byte{0xBE, 0xEF, 0xCA, 0xFE}, false, true},
		{"Request ten zeros", make([]byte, 10), true, false},
		{"Response ten zeros", make([]byte, 10), false, true},
		{"Request default packet size", make([]byte, defaultPacketSize), true, false},
		{"Response default packet size", make([]byte, defaultPacketSize), false, true},
		{"Request exceding packet size", make([]byte, defaultPacketSize+5), true, false},
		{"Response exceding packet size", make([]byte, defaultPacketSize+5), false, true},
		{"Request 3x default packet size", make([]byte, 3*defaultPacketSize), true, false},
		{"Response 3x default packet size", make([]byte, 3*defaultPacketSize), false, true},
		{"Request exceding 3x packet size", make([]byte, 3*defaultPacketSize+5), true, false},
		{"Response exceding 3x packet size", make([]byte, 3*defaultPacketSize+5), false, true},
	}

	// Prepare the connection
	p := NewTLSRelay()
	p.SetListen("127.0.0.1", 8080)
	p.SetRelay("127.0.0.1", 8081)

	lserver, lclient := net.Pipe()
	rserver, rclient := net.Pipe()

	go func() {
		defer lserver.Close()
		defer rclient.Close()
		e := p.relay(lserver, rclient)
		if e != nil {
			log.Println(e)
		}
	}()
	p.WaitToStart() // needed to start properly

	defer lclient.Close()
	defer rserver.Close()

	// Start with tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.content
			res := make([]byte, len(want)+1)

			if tt.isRequest {
				n, e := lclient.Write(want)
				if e != nil {
					t.Errorf("TLSRelay.relay() Client writing error = %v", e)
					return
				}

				n, e = rserver.Read(res)
				if e != nil {
					t.Errorf("TLSRelay.relay() Server reading error = %v", e)
					return
				}
				got := res[:n]

				if !reflect.DeepEqual(want, got) {
					if len(want) != len(got) {
						t.Errorf("TLSRelay.relay() Request len(got) = %v, len(want) = %v", len(got), len(want))
					} else {
						t.Errorf("TLSRelay.relay() Request got = %v, want = %v", got, want)
					}
				}
			}

			if tt.isResponse {
				n, e := rserver.Write(want)
				if e != nil {
					t.Errorf("TLSRelay.relay() Server writing error = %v", e)
					return
				}

				n, e = lclient.Read(res)
				if e != nil {
					t.Errorf("TLSRelay.relay() Client reading error = %v", e)
					return
				}
				got := res[:n]

				if !reflect.DeepEqual(want, got) {
					if len(want) != len(got) {
						t.Errorf("TLSRelay.relay() Response len(got) = %v, len(want) = %v", len(got), len(want))
					} else {
						t.Errorf("TLSRelay.relay() Response got = %v, want = %v", got, want)
					}
				}
			}
		})
	}
}

func TestTLSRelayFull(t *testing.T) {
	// Prepare tests
	tests := []struct {
		name       string
		request    string
		wantCode   int
		wantLength int64
	}{
		{"Wrong URI", "/wrong/", 404, 19},
		{"Response 10 bytes", "/echo/10", 200, 10},
		{"Response 100 bytes", "/echo/100", 200, 100},
		{"Response 1000 bytes", "/echo/1000", 200, 1000},
	}

	// Prepare listening server
	httpsPort := 11443
	var srv http.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/echo/", echoServer)
	srv.Addr = "127.0.0.1:" + strconv.Itoa(httpsPort)
	srv.Handler = mux
	defer srv.Close()
	go srv.ListenAndServeTLS("relay_test_cert/server.crt", "relay_test_cert/server.key")

	// Prepare client
	tlsconfig := &tls.Config{InsecureSkipVerify: true}
	transport := &http.Transport{TLSClientConfig: tlsconfig}
	client := &http.Client{Transport: transport}

	// Prepare TLSRelay
	p := NewTLSRelay()
	p.SetListen("127.0.0.1", 8080)
	p.SetRelay("127.0.0.1", httpsPort)
	defer p.Stop()
	go p.Start()

	p.WaitToStart() // needed to start properly

	// Start with tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, e := client.Get("https://127.0.0.1:8080" + tt.request)

			if e != nil {
				t.Errorf("TLSRelay FullTest error %v", e)
			}

			gotCode := res.StatusCode
			if gotCode != tt.wantCode {
				t.Errorf("TLSRelay FullTest response code got = %v, want = %v", gotCode, tt.wantCode)
			}

			gotLength := res.ContentLength
			if gotLength != tt.wantLength {
				t.Errorf("TLSRelay FullTest response length got = %v, want = %v", gotLength, tt.wantLength)
			}
		})
	}
}

func echoServer(w http.ResponseWriter, req *http.Request) {
	pieces := strings.Split(req.RequestURI, "/")

	if len(pieces) != 3 {
		return
	}

	n, e := strconv.Atoi(pieces[2])
	if e != nil {
		return
	}

	res := strings.Repeat("A", n)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(res))
}
