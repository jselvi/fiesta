package core

// optionType struct is...
type optionType struct {
	Param string
	Value string
	Descr string
}

// SetAttack function...
func (c *Core) SetAttack(a string) {
	c.initMap()
	for _, v := range Attacks {
		if v.cmd == a {
			c.attack = v.id
			return
		}
	}
}

// SetOption function...
func (c *Core) SetOption(option, value string) {
	c.initMap()
	c.option[option] = value
}

// UnsetOption function...
func (c *Core) UnsetOption(option string) {
	c.initMap()
	delete(c.option, option)
	c.setOptions(c.Options())
}

// Options function...
func (c *Core) Options() []optionType {
	var opt []optionType

	switch c.attack {
	case CRIME:
		opt = c.crimeOptions()
	case BREACH:
		opt = c.breachOptions()
	case TIME:
		opt = c.timeOptions()
	case XSSEARCH:
		opt = c.xssearchOptions()
	case XSTIME:
		opt = c.xstimeOptions()
	case FIESTA:
		opt = c.fiestaOptions()
	case DUMMY:
		opt = c.dummyOptions()
	default:
		opt = []optionType{}
	}

	c.setOptions(opt)
	return opt
}

func (c *Core) setOptions(opt []optionType) {
	for i, x := range opt {
		if len(x.Param) == 0 {
			continue
		}

		value, ok := c.option[x.Param]
		if ok {
			opt[i].Value = value
			continue
		}
		c.option[x.Param] = x.Value
	}
}

func (c *Core) defaultOptions() []optionType {
	return []optionType{
		{"SRVHOST", "0.0.0.0", "info"},
		{"SRVPORT", "8080", "info"},
		{"LOAD", "IMG", "Technique used to load guesses: IMG, IFRAME, OPEN, TAB"},
		{"PROXYHOST", "0.0.0.0", "info"},
		{"PROXYPORT", "10443", "info"},
		{"RELAYHOST", "127.0.0.1", "info"},
		{"RELAYPORT", "443", "info"},
		{"URL", "https://127.0.0.1/search/" + wildcard, "info"},
		{"WRONG", "^^^", "info"},
		{"CHARSET", "0123456789", "info"},
		{"SEARCH", "", "info"},
	}
}

func (c *Core) timingOptions() []optionType {
	removeList := map[string]bool{"PROXYHOST": true, "PROXYPORT": true, "RELAYHOST": true, "RELAYPORT": true}
	opts := c.defaultOptions()
	return c.removeOptions(opts, removeList)
}

func (c *Core) norelayOptions() []optionType {
	removeList := map[string]bool{"RELAYHOST": true, "RELAYPORT": true}
	opts := c.defaultOptions()
	return c.removeOptions(opts, removeList)
}

func (c *Core) removeOptions(opts []optionType, removeList map[string]bool) []optionType {
	var res []optionType
	for _, opt := range c.defaultOptions() {
		_, remove := removeList[opt.Param]
		if !remove {
			res = append(res, opt)
		}
	}
	return res
}
