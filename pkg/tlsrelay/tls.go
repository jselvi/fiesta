package tlsrelay

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

const (
	errorTLSChTrunked = "TLS Channel Full, Trunked"
)

func (p *TLSRelay) manage(client string, isRequest bool, data []byte) [][]byte {

	dataSlices := p.splitBytes(data)
	for _, ds := range dataSlices {
		// Info about the stream in JSON format
		json := p.bytesToJSON(client, isRequest, ds)
		p.writeInfo(json) // TODO: Do something with the error, not just ignore it
	}

	return dataSlices
}

func (p *TLSRelay) bytesToJSON(client string, isRequest bool, data []byte) string {
	pkg := gopacket.NewPacket(data, layers.LayerTypeTLS, p.gopacketOpt)
	if pkg.ErrorLayer() != nil {
		return p.rawBytesToJSON(client, isRequest, data)
	}

	tlsLayer := pkg.Layer(layers.LayerTypeTLS)
	if tlsLayer == nil {
		return p.rawBytesToJSON(client, isRequest, data)
	}

	tls, _ := tlsLayer.(*layers.TLS)
	if len(tls.AppData) == 1 {
		return p.appDataToJSON(client, isRequest, tls.AppData[0])
	}

	return p.rawBytesToJSON(client, isRequest, data)
}

func (p *TLSRelay) rawBytesToJSON(client string, isRequest bool, data []byte) string {
	type MsgStruct struct {
		Client    string
		IsRequest bool
		Type      string
		Length    int
	}
	var msg MsgStruct

	msg.Client = client
	msg.IsRequest = isRequest
	msg.Type = "Raw"
	msg.Length = len(data)

	jsonBytes, _ := json.Marshal(msg)
	return string(jsonBytes)
}

func (p *TLSRelay) appDataToJSON(client string, isRequest bool, r layers.TLSAppDataRecord) string {
	type MsgStruct struct {
		Client    string
		IsRequest bool
		Type      string
		Length    uint16
		Version   string
	}
	var msg MsgStruct

	msg.Client = client
	msg.IsRequest = isRequest

	msg.Type = "AppData"
	msg.Length = r.Length
	msg.Version = r.Version.String()

	jsonBytes, _ := json.Marshal(msg)
	return string(jsonBytes)
}

func (p *TLSRelay) tlsLength(data []byte) int {
	if len(data) < 5 {
		return 0
	}

	// Check TLS Record Content Type
	if data[0] < 20 || data[0] > 23 {
		return 0
	}

	// Check TLS Version
	if data[1] > 4 || data[2] > 5 {
		return 0
	}

	// Get Record length
	l := int(binary.BigEndian.Uint16(data[3:5]) + 5)

	// Too long
	if l > len(data) {
		//return 0 DEBUG
		l = len(data)
	}

	// Parse packet
	pkg := gopacket.NewPacket(data[:l], layers.LayerTypeTLS, p.gopacketOpt)
	if pkg.ErrorLayer() != nil {
		return 0
	}
	return l
}

func (p *TLSRelay) splitBytes(data []byte) [][]byte {
	var res [][]byte
	var rawAcc []byte
	var ini, end int

	l := len(data)
	for ini < l {
		end = ini + p.tlsLength(data[ini:])

		// No TLS starting from ini
		if ini == end {
			rawAcc = append(rawAcc, data[ini]) // accumulate the initial byte
			ini++
			continue
		}

		// TLS starting from ini

		// Add raw bytes first, if any
		if len(rawAcc) > 0 {
			res = append(res, rawAcc)
			rawAcc = make([]byte, 0)
		}

		// Add the TLS stream that we have found and continue
		res = append(res, data[ini:end])
		ini = end
	}

	if len(rawAcc) > 0 {
		res = append(res, rawAcc)
	}

	return res
}

func (p *TLSRelay) writeInfo(info string) error {
	var e error
	e = nil

	if len(p.tlsInfo) == p.tlsInfoChSize {
		p.Clean()
		e = errors.New(errorTLSChTrunked)
	}

	p.tlsInfo <- info
	return e
}

// Clean function removes all the TLS information from the Queue
func (p *TLSRelay) Clean() []string {
	var res []string
	for p.Len() > 0 {
		x := <-p.tlsInfo
		res = append(res, x)
	}
	return res
}

// Read function reads a piece of information from the Queue
func (p *TLSRelay) Read() string {
	return <-p.tlsInfo
}

// Len function return the size of the Queue
func (p *TLSRelay) Len() int {
	return len(p.tlsInfo)
}
