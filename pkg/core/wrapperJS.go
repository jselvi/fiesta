package core

import "errors"

// htmlContent function...
func (c *Core) htmlContent() string {
	return defaultHTML
}

// jsContent function...
func (c *Core) jsContent() string {
	switch c.attack {
	case CRIME:
		return c.crimeJsContent()
	case BREACH:
		return c.breachJsContent()
	case TIME:
		return c.timeJsContent()
	case XSSEARCH:
		return c.xssearchJsContent()
	case XSTIME:
		return c.xstimeJsContent()
	case FIESTA:
		return c.fiestaJsContent()
	case DUMMY:
		return c.dummyJsContent()
	default:
		return ""
	}
}

func (c *Core) defaultJS() string {

	x, ok := c.option["LOAD"]
	if !ok {
		return ""
	}
	l, e := StringToLoad(x)
	if e != nil {
		return ""
	}

	var s string
	switch l {
	case IMG:
		s += loadImg
	case IFRAME:
		s += loadIFrame
	case OPEN:
		s += loadOpen
	case TAB:
		s += loadTab
	}
	s += loadFromWS
	return s
}

func (c *Core) timingJS() string {
	s := c.defaultJS() + measureTiming
	return s
}

// getInfo function...
func (c *Core) getInfo() ([]string, error) {
	switch c.attack {
	case CRIME:
		return c.crimeGetInfo()
	case BREACH:
		return c.breachGetInfo()
	case TIME:
		return c.timeGetInfo()
	case XSSEARCH:
		return c.xssearchGetInfo()
	case XSTIME:
		return c.xstimeGetInfo()
	case FIESTA:
		return c.fiestaGetInfo()
	case DUMMY:
		return c.dummyGetInfo()
	default:
		return []string{}, errors.New("getInfo not implemented for the given attack")
	}
}
