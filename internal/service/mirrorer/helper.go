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

func FindPath(inp []byte) []string {
	//re := regexp.MustCompile(`(url\(['"]\.?(\S+)['"]\))|(href=\'(\w+\:\/\/\S+)\') |(src=\'(\S+)\')`)
	//re := regexp.MustCompile(`(url\(['"]\.?(?P<url0>\S+)['"]\))|(href=\'(?P<url1>\w+\:\/\/\S+)\') |(src=\'(?P<url2>\S+)\')`)
	//re := regexp.MustCompile(`(url\(['"]\.?(?P<url0>\S+)['"]\))|(href=\'(?P<url1>\w+\:\/\/\S+)\') |(src=\'(?P<url2>\S+)\') |(src=['"](?P<url3>\S+)['"])`)
	re := regexp.MustCompile(`(url\(['"]\.?(?P<url0>\S+)['"]\))|(href=\'(?P<url1>\w+\:\/\/\S+)\')|(src=\'(?P<url2>\S+)\')|(src=['"](?P<url3>\S+)['"])|(href=['"](?P<url4>[\w/\.]+)['"])`)

	resRegex := re.FindAllSubmatch(inp, -1)
	//fmt.Println(len(resRegex))
	var res []string
	for _, v := range resRegex {
		//fmt.Println(re.SubexpIndex("url0"), re.SubexpIndex("url1"), re.SubexpIndex("url2"), re.SubexpIndex("url3"), re.SubexpIndex("url4"))
		switch {
		case len(v) > 10 && len(v[10]) != 0:
			res = append(res, string(v[10]))
		case len(v) > 8 && len(v[8]) != 0:
			fmt.Println(string(v[8]))
			res = append(res, string(v[8]))
		case len(v) >= 7 && len(v[6]) != 0:
			res = append(res, string(v[6]))
		case len(v) >= 5 && len(v[4]) != 0:
			res = append(res, string(v[4]))
		case len(v) >= 3 && len(v[2]) != 0:
			res = append(res, string(v[2]))
		}

	}
	return res
}

func ContainsProto(url string) bool {
	re := regexp.MustCompile(`^(\D{3,7}\:\/\/)`)
	return re.MatchString(url)
}
