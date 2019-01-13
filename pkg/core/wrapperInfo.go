package core

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

func (c *Core) getProxyInfo() []string {
	// Get Information
	var tlsInfo []string
	startGuess := time.Now()
	startConn := time.Now()
	for {
		// Calculate elapsed time
		elapsedGuess := time.Since(startGuess)
		elapsedConn := time.Since(startConn)

		// Break if timeout
		if elapsedGuess > c.timeout {
			break
		}
		if elapsedConn > c.timeoutConn {
			break
		}

		// Wait until there are info to process
		if c.relay.Len() == 0 {
			continue
		}

		// Append results
		s := c.relay.Read()
		tlsInfo = append(tlsInfo, s)

		// Recalculate connection start
		startConn = time.Now()
	}

	// Sort information
	tlsSorted, _ := c.sortInfo(tlsInfo)

	// Compact information
	tlsCompacted, _ := c.compactAppData(tlsSorted)

	return tlsCompacted
}

type proxyInfo struct {
	Client    string
	IsRequest bool
	Type      string
	Length    int
}

func (c *Core) sortInfo(info []string) ([]string, error) {
	infoMap := make(map[string][]string)
	var infoMapOrder []string
	idx := make(map[string]int)

	for _, x := range info {
		// Decode JSON
		b := []byte(x)
		var m proxyInfo
		e := json.Unmarshal(b, &m)
		if e != nil {
			return []string{}, e
		}

		// Create index if doesn't exist
		_, ok := idx[m.Client]
		if !ok {
			idx[m.Client] = 0
		}

		// Calculate where to store
		if m.IsRequest {
			idx[m.Client]++
		}
		key := m.Client + "-" + strconv.Itoa(idx[m.Client])

		// Store
		_, ok = infoMap[key]
		if !ok {
			infoMap[key] = []string{}
			infoMapOrder = append(infoMapOrder, key)
		}
		infoMap[key] = append(infoMap[key], x)
	}

	// Merge in a single slice
	var sorted []string
	for _, key := range infoMapOrder {
		x := infoMap[key]
		sorted = append(sorted, x...)
	}

	return sorted, nil
}

func (c *Core) compactAppData(info []string) ([]string, error) {
	var compacted []proxyInfo
	i := -1
	lastWasRequest := true

	for _, x := range info {
		// Decode JSON
		b := []byte(x)
		var m proxyInfo
		e := json.Unmarshal(b, &m)
		if e != nil {
			return []string{}, e
		}

		if m.Type != "AppData" {
			compacted = append(compacted, m)
			continue
		}

		if m.IsRequest {
			compacted = append(compacted, m)
			lastWasRequest = true
			continue
		}

		if lastWasRequest {
			lastWasRequest = false
			i = len(compacted)
			compacted = append(compacted, m)
			continue
		}

		compacted[i].Length += m.Length
	}

	var compactedStr []string
	for _, x := range compacted {
		b, _ := json.Marshal(x)
		s := string(b)
		compactedStr = append(compactedStr, s)
	}

	return compactedStr, nil
}

func (c *Core) getWebSocketInfo() []string {
	return []string{}
}

// getOracle function...
func (c *Core) getOracle(s string) (int, error) {
	switch c.attack {
	case CRIME:
		return c.crimeGetOracle(s)
	case BREACH:
		return c.breachGetOracle(s)
	case TIME:
		return c.timeGetOracle(s)
	case XSSEARCH:
		return c.xssearchGetOracle(s)
	case XSTIME:
		return c.xstimeGetOracle(s)
	case FIESTA:
		return c.fiestaGetOracle(s)
	case DUMMY:
		return c.dummyGetOracle(s)
	default:
		return 0, errors.New("getOracle not implemented for the given attack")
	}
}
