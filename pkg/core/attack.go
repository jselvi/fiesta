package core

const (
	CRIME    = 1
	BREACH   = 2
	TIME     = 3
	XSSEARCH = 4
	XSTIME   = 5
	FIESTA   = 6
	DUMMY    = 7
)

var Attacks = AttacksBatch{
	//{CRIME, "crime", "CRIME"},
	//{BREACH, "breach", "BREACH"},
	//{TIME, "time", "TIME"},
	{XSSEARCH, "xssearch", "XSSEARCH"},
	//{XSTIME, "xstime", "XSTIME"},
	{FIESTA, "fiesta", "FIESTA"},
	//{DUMMY, "dummy", "DUMMY"},
}

type AttacksBatch []AttackInfo
type AttackInfo struct {
	id   uint8
	cmd  string
	desc string
}

func (a AttackInfo) String() string {
	return a.desc
}

func (a AttackInfo) Cmd() string {
	return a.cmd
}

func (a AttacksBatch) DumpCmd() []string {
	var res []string
	for _, v := range a {
		res = append(res, v.cmd)
	}
	return res
}
