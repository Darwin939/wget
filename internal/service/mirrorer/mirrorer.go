package mirrorer

import (
	"fmt"
	"strings"
)

type Mirrorer struct {
	Excluded      []string
	Reject        []string
	ExcludedRegex string
	RejectRegex   string
	URL           string
}

func NewMirrorer(url string, excluded, reject []string) *Mirrorer {
	return &Mirrorer{
		Excluded: excluded,
		Reject:   reject,
		URL:      url,
	}
}

func (m *Mirrorer) CreateMirror() {
	m.initRegex()

}

func (m *Mirrorer) initRegex() {
	//TODO find optimal way for creating regex
	m.ExcludedRegex = regexBuilder(`\.`, m.Excluded)
	m.RejectRegex = regexBuilder(`\/`, m.Reject)
}

func regexBuilder(delim string, params []string) string {
	if len(params) == 0 {
		return ""
	}
	res := strings.Builder{}
	res.WriteString("[^(")
	if len(params[0]) != 0 {
		res.WriteString(fmt.Sprintf("(%s%s)", delim, params[0]))
	}
	for _, v := range params[1:] {
		if len(v) != 0 {
			res.WriteString(fmt.Sprintf("|(%s%s)", delim, v))
		}
	}
	res.WriteString(")]")
	return res.String()
}
