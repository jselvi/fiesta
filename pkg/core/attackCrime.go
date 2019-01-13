package core

import (
	"encoding/json"
	"errors"
)

func (c *Core) crimeOptions() []optionType {
	return c.defaultOptions()
}

func (c *Core) crimeJsContent() string {
	return c.defaultJS()
}

func (c *Core) crimeGetInfo() ([]string, error) {
	return c.getProxyInfo(), nil
}

type crimeTLS struct {
	Client    string
	IsRequest bool
	Type      string
	Length    int
}

func (c *Core) crimeGetOracle(s string) (int, error) {
	b := []byte(s)
	var m crimeTLS
	e := json.Unmarshal(b, &m)
	if e != nil {
		return 0, e
	}

	e = nil
	if m.Type != "AppData" {
		e = errors.New("packet is not appdata")
	}
	if !m.IsRequest {
		e = errors.New("packet is a response")
	}

	return m.Length, e
}

func (c *Core) crimeGuess(g string) (int, error) {
	return 0, nil
}

func (c *Core) crimeMakeDecision(right, wrong int, results map[string]int) ([]string, error) {
	var ok []string
	for key, value := range results {
		if value == right {
			ok = append(ok, key)
		}
	}
	return ok, nil
}
