package tlsrelay

import (
	"reflect"
	"testing"
)

func TestTLSRelayNew(t *testing.T) {
	p := NewTLSRelay()
	ptype := reflect.TypeOf(p)
	if ptype.String() != "tlsrelay.TLSRelay" {
		t.Errorf("TLSRelay.NewTLSRelay() does not return a tlsrelay.TLSRelay type")
	}
}

func TestTLSRelayStart(t *testing.T) {
	tests := []struct {
		name       string
		listenHost string
		listenPort int
		relayHost  string
		relayPort  int
		wantErr    bool
	}{
		{"Wrong listenHost", "327.0.0.1", 8080, "127.0.0.1", 8081, true},
		{"Wrong listenPort", "127.0.0.1", -10, "127.0.0.1", 8081, true},
		{"Wrong relayHost", "127.0.0.1", 8080, "327.0.0.1", 8081, true},
		{"Wrong relayPort", "127.0.0.1", 8080, "127.0.0.1", -10, true},
		{"No Listen", defaultHost, defaultPort, "127.0.0.1", 8081, false},
		{"No Relay", "127.0.0.1", 8080, "", 0, false},
		{"Right Relay", "127.0.0.1", 8080, "127.0.0.1", 8081, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTLSRelay()
			p.listenHost = tt.listenHost
			p.listenPort = tt.listenPort
			p.relayHost = tt.relayHost
			p.relayPort = tt.relayPort

			go func() {
				p.WaitToStart() // needed to start properly
				p.Stop()
			}()
			e := p.Start()
			got := (e != nil)

			if got != tt.wantErr {
				t.Errorf("TLSRelay.Start() error = %v(%v), wantErr %v", got, e, tt.wantErr)
				return
			}
		})
	}
}

func TestTLSRelayStop(t *testing.T) {
	// Stop a server that was not even created
	var p1 TLSRelay
	e1 := p1.Stop()
	if e1 == nil {
		t.Errorf("TLSRelay.Stop() error = nil, wanted 'error when TLSRelay is not created'")
	}

	// Create a TLSRelay
	var p = NewTLSRelay()
	go p.Start()
	p.WaitToStart() // needed to start properly
	e1 = p.Stop()
	if e1 != nil {
		t.Errorf("TLSRelay.Stop() error %v, wanted nil", e1)
	}

}

func TestTLSRelaySetListen(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		wantErr bool
	}{
		{"Right IPv4", "127.0.0.1", 80, false},
		{"Right IPv6", "2001:db8:a0b:12f0::1", 80, false},
		{"Right Hostname", "www.pentester.es", 80, false},
		{"Wrong IPv4", "327.0.0.1", 80, true},
		{"Wrong IPv6", "2001:db8:a0b:12f0::1:5:4:3:2:1", 80, true},
		{"Wrong Hostname", "doesnotexist.pentester.es", 80, true},
		{"Zero Port", "127.0.0.1", 0, true},
		{"Negative Port", "127.0.0.1", -80, true},
		{"Too High Port", "127.0.0.1", 65536, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTLSRelay()
			if err := p.SetListen(tt.host, tt.port); (err != nil) != tt.wantErr {
				t.Errorf("TLSRelay.SetListen() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTLSRelaySetRelay(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		wantErr bool
	}{
		{"Right IPv4", "127.0.0.1", 80, false},
		{"Right IPv6", "2001:db8:a0b:12f0::1", 80, false},
		{"Right Hostname", "www.pentester.es", 80, false},
		{"Wrong IPv4", "327.0.0.1", 80, true},
		{"Wrong IPv6", "2001:db8:a0b:12f0::1:5:4:3:2:1", 80, true},
		{"Wrong Hostname", "doesnotexist.pentester.es", 80, true},
		{"Zero Port", "127.0.0.1", 0, true},
		{"Negative Port", "127.0.0.1", -80, true},
		{"Too High Port", "127.0.0.1", 65536, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTLSRelay()
			if err := p.SetRelay(tt.host, tt.port); (err != nil) != tt.wantErr {
				t.Errorf("TLSRelay.SetRelay() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTLSRelayCheckHostPort(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		wantErr bool
	}{
		{"Right IPv4", "127.0.0.1", 80, false},
		{"Right IPv6", "2001:db8:a0b:12f0::1", 80, false},
		{"Right Hostname", "www.pentester.es", 80, false},
		{"Wrong IPv4", "327.0.0.1", 80, true},
		{"Wrong IPv6", "2001:db8:a0b:12f0::1:5:4:3:2:1", 80, true},
		{"Wrong Hostname", "doesnotexist.pentester.es", 80, true},
		{"Zero Port", "127.0.0.1", 0, true},
		{"Negative Port", "127.0.0.1", -80, true},
		{"Too High Port", "127.0.0.1", 65536, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewTLSRelay()
			err := p.checkHostPort(tt.host, tt.port)
			got := (err != nil)
			if got != tt.wantErr {
				t.Errorf("TLSRelay.checkHostPort() error = %v, wantErr %v", got, tt.wantErr)
			}
		})
	}
}

func TestTLSRelayResetRelay(t *testing.T) {
	t.Run("Reset", func(t *testing.T) {
		p := NewTLSRelay()
		p.relayHost = "thisisarandomhostthatdoesnotexist"
		p.relayPort = -10
		p.ResetRelay()

		if (p.relayHost != defaultHost) || (p.relayPort != defaultPort) {
			t.Errorf("TLSRelay.ResetRelay() error = true, wantErr false")
		}
	})
}

func TestTLSRelayDemultiplex(t *testing.T) {
	p := NewTLSRelay()

	if p.sendDelay != defaultSendDelay {
		t.Errorf("TLSRelay Default sendDelay got = %d, want %d", p.sendDelay, defaultSendDelay)
	}

	p.Demultiplex(true)
	if p.sendDelay != maxSendDelay {
		t.Errorf("TLSRelay.Demultiplex(true) got = %d, want %d", p.sendDelay, maxSendDelay)
	}

	p.Demultiplex(false)
	if p.sendDelay != defaultSendDelay {
		t.Errorf("TLSRelay.Demultiplex(true) got = %d, want %d", p.sendDelay, defaultSendDelay)
	}
}
