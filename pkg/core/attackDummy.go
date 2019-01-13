package core

func (c *Core) dummyOptions() []optionType {
	return c.defaultOptions()
}

func (c *Core) dummyJsContent() string {
	return c.defaultJS()
}

func (c *Core) dummyGetInfo() ([]string, error) {
	return []string{}, nil
}

func (c *Core) dummyGetOracle(s string) (int, error) {
	return 0, nil
}

func (c *Core) dummyGuess(g string) (int, error) {
	return 0, nil
}

func (c *Core) dummyMakeDecision(right, wrong int, results map[string]int) ([]string, error) {
	return []string{}, nil
}
