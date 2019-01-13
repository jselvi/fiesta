package core

import (
	"encoding/json"
	"errors"
)

func (c *Core) fiestaOptions() []optionType {
	return c.norelayOptions()
}

func (c *Core) fiestaJsContent() string {
	return c.defaultJS()
}

func (c *Core) fiestaGetInfo() ([]string, error) {
	return c.getProxyInfo(), nil
}

type fiestaTLS struct {
	Client    string
	IsRequest bool
	Type      string
	Length    int
}

func (c *Core) fiestaGetOracle(s string) (int, error) {
	b := []byte(s)
	var m fiestaTLS
	e := json.Unmarshal(b, &m)
	if e != nil {
		return 0, e
	}

	e = nil
	if !m.IsRequest {
		e = errors.New("packet is a response")
	}

	return 1, e
}

func (c *Core) fiestaGuess(g string) (int, error) {
	return 0, nil
}

func (c *Core) fiestaMakeDecision(right, wrong int, results map[string]int) ([]string, error) {
	var ok []string
	for key, value := range results {
		//if value < wrong-1 || value > wrong+1 {
		if value != wrong {
			ok = append(ok, key)
		}
	}
	return ok, nil
}

func (c *Core) fiestaOracleToInt(oracles []int) (int, error) {
	return len(oracles), nil
}
