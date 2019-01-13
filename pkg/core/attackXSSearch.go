package core

import (
	"encoding/json"
	"errors"
)

func (c *Core) xssearchOptions() []optionType {
	return c.defaultOptions()
}

func (c *Core) xssearchJsContent() string {
	return c.defaultJS()
}

func (c *Core) xssearchGetInfo() ([]string, error) {
	return c.getProxyInfo(), nil
}

type xssearchTLS struct {
	Client    string
	IsRequest bool
	Type      string
	Length    int
}

func (c *Core) xssearchGetOracle(s string) (int, error) {
	b := []byte(s)
	var m xssearchTLS
	e := json.Unmarshal(b, &m)
	if e != nil {
		return 0, e
	}

	e = nil
	if m.Type != "AppData" {
		e = errors.New("packet is not appdata")
	}
	if m.IsRequest {
		e = errors.New("packet is a request")
	}

	return m.Length, e
}

func (c *Core) xssearchGuess(g string) (int, error) {
	return 0, nil
}

func (c *Core) xssearchMakeDecision(right, wrong int, results map[string]int) ([]string, error) {
	var ok []string
	for key, value := range results {
		if value > wrong {
			ok = append(ok, key)
		}
	}
	return ok, nil
}
