package mirrorer

import (
	"fmt"
	"regexp"
	"strings"
)

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

func convertToPath(url string) string {
	re := regexp.MustCompile(`^\w+:\/\/`)
	return re.ReplaceAllString(url, "")
}

func validateURL(url string) string {
	re := regexp.MustCompile(`^(http)|(https)|(ftp)\:\/\/`)
	if re.MatchString(url) {
		return url
	}
	re = regexp.MustCompile(`^\w{1,6}:\/\/`)
	if !re.MatchString(url) {
		return fmt.Sprintf("http://%s", url)
	}
	return url
}
