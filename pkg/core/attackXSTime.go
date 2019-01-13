package core

func (c *Core) xstimeOptions() []optionType {
	return c.timingOptions()
}

func (c *Core) xstimeJsContent() string {
	return c.timingJS()
}

func (c *Core) xstimeGetInfo() ([]string, error) {
	return c.getWebSocketInfo(), nil
}

func (c *Core) xstimeGetOracle(s string) (int, error) {
	return 0, nil
}

func (c *Core) xstimeGuess(g string) (int, error) {
	return 0, nil
}

func (c *Core) xstimeMakeDecision(right, wrong int, results map[string]int) ([]string, error) {
	var ok []string
	for key, value := range results {
		diff := value - wrong
		x := diff * 100 / value
		if x > 10 {
			ok = append(ok, key)
		}
	}
	return ok, nil
}
