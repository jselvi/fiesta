package core

func (c *Core) initMap() {
	c.mtx.Lock()
	if c.option == nil {
		c.option = make(map[string]string)
	}
	c.mtx.Unlock()
}

func (c *Core) sendIfOpen(ch chan string, msg string) bool {
	if !c.running {
		return false
	}

	ch <- msg
	return true
}

func (c *Core) pop(slice []string) ([]string, string) {
	if len(slice) == 0 {
		return []string{}, ""
	}
	return slice[:len(slice)-1], slice[len(slice)-1]
}

func (c *Core) push(slice []string, value string) []string {
	return append(slice, value)
}
