package core

import (
	"errors"
	"strings"
)

// LoadTechnique is the technique used to load a guess
type LoadTechnique uint8

// Three techniques are supported: img, iframe and open
const (
	IMG    = LoadTechnique(0)
	IFRAME = LoadTechnique(1)
	OPEN   = LoadTechnique(2)
	TAB    = LoadTechnique(3)
)

func (l LoadTechnique) String() string {
	switch l {
	case IMG:
		return "IMG"
	case IFRAME:
		return "IFRAME"
	case OPEN:
		return "OPEN"
	case TAB:
		return "TAB"
	default:
		return ""
	}
}

// ListLoadTechniques function returns all the possible values for Load Techniques
func ListLoadTechniques() []string {
	return []string{IMG.String(), IFRAME.String(), OPEN.String(), TAB.String()}
}

// StringToLoad function convert a string to the apropiate number
func StringToLoad(s string) (LoadTechnique, error) {
	for i, x := range ListLoadTechniques() {
		if strings.ToUpper(s) == x {
			return LoadTechnique(i), nil
		}
	}
	return LoadTechnique(0), errors.New("Load Technique Does Not Exist")
}
