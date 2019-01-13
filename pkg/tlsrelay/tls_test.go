package tlsrelay

import (
	"reflect"
	"testing"
)

var testClient = "192.168.1.1:12345"

var testDoubleAppData = []byte{
	0x17, 0x03, 0x01, 0x00, 0x20, 0x77, 0x3a, 0x94, 0x7d, 0xb4, 0x47, 0x4a, 0x1d, 0xd4, 0x6c, 0x5a,
	0x69, 0x74, 0x03, 0x93, 0x32, 0xca, 0x54, 0x5e, 0xa5, 0x81, 0x99, 0x6a, 0x73, 0x66, 0xbf, 0x06,
	0xa0, 0xdc, 0x6a, 0x9c, 0xb1, 0x17, 0x03, 0x01, 0x00, 0x20, 0x44, 0x64, 0xc8, 0xc2, 0x5a, 0xfc,
	0x4a, 0x82, 0xdd, 0x53, 0x6d, 0x30, 0x82, 0x4d, 0x35, 0x22, 0xf1, 0x5f, 0x3b, 0x96, 0x66, 0x79,
	0x61, 0x9f, 0x51, 0x93, 0x1b, 0xbf, 0x53, 0x3b, 0xf8, 0x26,
}
var testDoubleAppDataIsRequest = true
var testDoubleAppDataOneJSON = "{\"Client\":\"192.168.1.1:12345\",\"IsRequest\":true,\"Type\":\"AppData\",\"Length\":32,\"Version\":\"TLS 1.0\"}"
var testDoubleAppDataTwoJSON = "{\"Client\":\"192.168.1.1:12345\",\"IsRequest\":true,\"Type\":\"AppData\",\"Length\":32,\"Version\":\"TLS 1.0\"}"

var testNoisyPacket = []byte{
	0xbe, 0xef, 0xca, 0xfe, 0xbe, 0xef, 0xca, 0xfe, 0xbe, 0xef, 0xca, 0xfe, 0xbe, 0xef, 0xca, 0xfe,
	0x17, 0x03, 0x01, 0x00, 0x20, 0x77, 0x3a, 0x94, 0x7d, 0xb4, 0x47, 0x4a, 0x1d, 0xd4, 0x6c, 0x5a,
	0x69, 0x74, 0x03, 0x93, 0x32, 0xca, 0x54, 0x5e, 0xa5, 0x81, 0x99, 0x6a, 0x73, 0x66, 0xbf, 0x06,
	0xa0, 0xdc, 0x6a, 0x9c, 0xb1, 0xff, 0x17, 0x03, 0x01, 0x00, 0x20, 0x44, 0x64, 0xc8, 0xc2, 0x5a,
	0xfc, 0x4a, 0x82, 0xdd, 0x53, 0x6d, 0x30, 0x82, 0x4d, 0x35, 0x22, 0xf1, 0x5f, 0x3b, 0x96, 0x66,
	0x79, 0x61, 0x9f, 0x51, 0x93, 0x1b, 0xbf, 0x53, 0x3b, 0xf8, 0x26, 0xff,
}
var testNoisyPacketPieces = []int{16, 37, 1, 37, 1}
var testNoisyIsRequest = true
var testNoisyOneJSON = "{\"Client\":\"192.168.1.1:12345\",\"IsRequest\":true,\"Type\":\"Raw\",\"Length\":16}"
var testNoisyTwoJSON = "{\"Client\":\"192.168.1.1:12345\",\"IsRequest\":true,\"Type\":\"AppData\",\"Length\":32,\"Version\":\"TLS 1.0\"}"
var testNoisyThreeJSON = "{\"Client\":\"192.168.1.1:12345\",\"IsRequest\":true,\"Type\":\"Raw\",\"Length\":1}"
var testNoisyFourJSON = testNoisyTwoJSON
var testNoisyFiveJSON = testNoisyThreeJSON

var testAppData = []byte{
	0x17, 0x03, 0x01, 0x00, 0x20, 0x77, 0x3a, 0x94, 0x7d, 0xb4, 0x47, 0x4a, 0x1d, 0xd4, 0x6c, 0x5a,
	0x69, 0x74, 0x03, 0x93, 0x32, 0xca, 0x54, 0x5e, 0xa5, 0x81, 0x99, 0x6a, 0x73, 0x66, 0xbf, 0x06,
	0xa0, 0xdc, 0x6a, 0x9c, 0xb1,
}
var testAppDataIsRequest = true
var testAppDataJSON = "{\"Client\":\"192.168.1.1:12345\",\"IsRequest\":true,\"Type\":\"AppData\",\"Length\":32,\"Version\":\"TLS 1.0\"}"

var testRaw = []byte{
	0xbe, 0xef, 0xca, 0xfe, 0xbe, 0xef, 0xca, 0xfe, 0xbe, 0xef, 0xca, 0xfe, 0xbe, 0xef, 0xca, 0xfe,
	0xbe, 0xef, 0xca, 0xfe, 0xbe, 0xef, 0xca, 0xfe,
}
var testRawIsRequest = false
var testRawJSON = "{\"Client\":\"192.168.1.1:12345\",\"IsRequest\":false,\"Type\":\"Raw\",\"Length\":24}"

func join(s [][]byte) []byte {
	var res []byte
	for _, x := range s {
		res = append(res, x...)
	}
	return res
}

func TestTLSRelayManage(t *testing.T) {
	// Prepare tests
	tests := []struct {
		name      string
		data      []byte
		isRequest bool
		want      []string
	}{
		{"Single AppData Request", testAppData, testDoubleAppDataIsRequest, []string{testAppDataJSON}},
		{"Double AppData Request", testDoubleAppData, testDoubleAppDataIsRequest, []string{testDoubleAppDataOneJSON, testDoubleAppDataOneJSON}},
		{"Raw Response", testRaw, testRawIsRequest, []string{testRawJSON}},
		{"Mixed Request", testNoisyPacket, testNoisyIsRequest, []string{testNoisyOneJSON, testNoisyTwoJSON, testNoisyThreeJSON, testNoisyFourJSON, testNoisyFiveJSON}},
	}

	// Prepare the connection
	p := NewTLSRelay()
	p.SetListen("127.0.0.1", 8080)
	p.SetRelay("127.0.0.1", 8081)

	// Start with tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Join Output. Needed to compare with tt.data later
			gotOutput := join(p.manage(testClient, tt.isRequest, tt.data))

			// Relaying properly?
			wantOutput := tt.data
			if !reflect.DeepEqual(wantOutput, gotOutput) {
				t.Errorf("%s Output = %v, want %v", tt.name, gotOutput, wantOutput)
			}

			// TLS Info
			got := p.Clean()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s JSON = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestTLSwriteInfo(t *testing.T) {
	// Prepare tests
	tests := []struct {
		name      string
		howMany   int
		wantError bool
	}{
		{"Write a message, then read", 1, false},
		{"Write max messages, then read", defaultTLSInfoChBufferSize, false},
		{"Write too many messages, error expected", defaultTLSInfoChBufferSize + 1, true},
	}

	// Prepare the connection
	p := NewTLSRelay()

	// Start with tests
	want := "TEST STRING"
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p.Clean()
			gotError := false
			for i := 0; i < tt.howMany; i++ {
				e := p.writeInfo(want)
				gotError = (gotError || (e != nil))
			}

			if gotError != tt.wantError {
				t.Errorf("%s Error = %v, want %v", tt.name, gotError, tt.wantError)
			}

			l := p.Len()
			for i := 0; i < l; i++ {
				got := p.Read()
				if got != want {
					t.Errorf("%s Output = %s, want %s", tt.name, got, want)
				}
			}
		})
	}
}

func TestTLSsplitBytes(t *testing.T) {
	// Prepare tests
	tests := []struct {
		name string
		data []byte
		want []int
	}{
		{"Split Single AppData", testAppData, []int{37}},
		{"Split Double AppData", testDoubleAppData, []int{37, 37}},
		{"Split Double AppData with Noise", testNoisyPacket, testNoisyPacketPieces},
	}

	// Prepare the connection
	p := NewTLSRelay()

	// Start with tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got []int
			gotSplit := p.splitBytes(tt.data)
			for _, v := range gotSplit {
				got = append(got, len(v))
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s pieces = %v, want %v", tt.name, got, tt.want)
			}

			gotOutput := join(gotSplit)
			if !reflect.DeepEqual(gotOutput, tt.data) {
				t.Errorf("%s Output = %v, want %v", tt.name, gotOutput, tt.data)
			}
		})
	}
}

func TestTLSbytesToJSON(t *testing.T) {
	// Prepare tests
	tests := []struct {
		name      string
		data      []byte
		isRequest bool
		want      string
	}{
		{"AppData Record", testAppData, testAppDataIsRequest, testAppDataJSON},
		{"Raw bytes", testRaw, testRawIsRequest, testRawJSON},
	}

	// Prepare the connection
	p := NewTLSRelay()

	// Start with tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.bytesToJSON(testClient, tt.isRequest, tt.data)
			if got != tt.want {
				t.Errorf("%s = %s, want %s", tt.name, got, tt.want)
			}
		})
	}
}
