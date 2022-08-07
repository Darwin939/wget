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
	re := regexp.MustCompile(`^(https?)|(ftp)\:\/\/`)
	if re.MatchString(url) {
		return url
	}
	re = regexp.MustCompile(`^\w{1,6}:\/\/`)
	if !re.MatchString(url) {
		return fmt.Sprintf("http://%s", url)
	}
	return url
}

func FindPath(inp string) []string {
	re := regexp.MustCompile(`(url\(['"]\.?(?P<url0>\/\S+)['"]\))`)
	resRegex := re.FindAllStringSubmatch(inp, -1)
	var res []string

	for _, v := range resRegex {
		if len(v) > 2 && v[2] != "" {
			res = append(res, v[2])
		}
	}
	return res
}

func isLocalPath(url string) (isLocal bool, path string) {
	re := regexp.MustCompile(`^(https?\:)?\/\/`)

	if re.MatchString(url) {
		return false, ""
	}
	re = regexp.MustCompile(`\.?\/?(?P<url>\w[\w\-\.]+(\/[\w\-\.]+)*)`)
	res := re.FindStringSubmatch(url)
	if len(res) > 1 && res[1] != "" {
		return true, res[1]
	}
	return true, url
}

func ContainsProto(url string) bool {
	re := regexp.MustCompile(`^(\D{3,7}\:\/\/)`)
	return re.MatchString(url)
}
