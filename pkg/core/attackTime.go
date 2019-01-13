package core

import "math"

func (c *Core) timeOptions() []optionType {
	return c.timingOptions()
}

func (c *Core) timeJsContent() string {
	return c.timingJS()
}

func (c *Core) timeGetInfo() ([]string, error) {
	return c.getWebSocketInfo(), nil
}

func (c *Core) timeGetOracle(s string) (int, error) {
	return 0, nil
}

func (c *Core) timeGuess(g string) (int, error) {
	return 0, nil
}

func (c *Core) timeMakeDecision(right, wrong int, results map[string]int) ([]string, error) {
	var ok []string
	for key, value := range results {
		diff := int(math.Abs(float64(right) - float64(value)))
		x := diff * 100 / value
		if x < 10 {
			ok = append(ok, key)
		}
	}
	return ok, nil
}
