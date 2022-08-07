package parser

import (
	"os"
	"path/filepath"
	"strings"
)

type Flags struct {
	Url          string // http https ftp
	IsFilename   bool   // flag -O=
	IsPathPassed bool   // flag -P=
	IsSpeedLimit bool   // flag --rate-limit=
	IsSaveFrom   bool   // flag -i=
	IsMirror     bool   // flag --mirror
	IsReject     bool   // flag --reject= -R=
	IsExclude    bool   // flag --exclude= -X=
	IsBackground bool   // flag -b
	Filename     string
	Path         string
	SpeedLimit   string
	SaveFrom     string
	Reject       string
	Exclude      string
}

func NewFlags(args []string) *Flags {
	return &Flags{
		Url:          get_Url(args),
		IsFilename:   get_Filename(args) != "",
		IsPathPassed: get_Path(args) != "",
		IsSpeedLimit: get_SpeedLimit(args) != "",
		IsSaveFrom:   get_SaveFrom(args) != "",
		IsMirror:     get_Mirror(args),
		IsReject:     get_Reject(args) != "",
		IsExclude:    get_Exclude(args) != "",
		Filename:     get_Filename(args),
		Path:         get_Path(args),
		SpeedLimit:   get_SpeedLimit(args),
		SaveFrom:     get_SaveFrom(args),
		Reject:       get_Reject(args),
		Exclude:      get_Exclude(args),
		IsBackground: getBackground(args) != "",
	}
}

func get_Url(args []string) string {
	for _, v := range args {
		if strings.HasPrefix(v, "ftp://") || strings.HasPrefix(v, "http://") || strings.HasPrefix(v, "https://") {
			return v
		}
	}

	return ""
}

func get_Filename(args []string) string {
	for _, v := range args {
		if strings.HasPrefix(v, "-O=") {
			return v[3:]
		}
	}

	return ""
}

func get_Path(args []string) string {
	ans := ""
	for _, v := range args {
		if strings.HasPrefix(v, "-P=") {
			ans += v[3:]
		}
	}

	if ans != "" && ans[0] == '~' {
		dirname, _ := os.UserHomeDir()
		ans = filepath.Join(dirname, ans[1:])
	}

	return ans
}

func get_SpeedLimit(args []string) string {
	for _, v := range args {
		if strings.HasPrefix(v, "--rate-limit=") {
			return v[13:]
		}
	}
	return ""
}

func get_SaveFrom(args []string) string {
	for _, v := range args {
		if strings.HasPrefix(v, "-i=") {
			return v[3:]
		}
	}

	return ""
}

func get_Mirror(args []string) bool {
	for _, v := range args {
		if v == "--mirror" {
			return true
		}
	}

	return false
}

func get_Reject(args []string) string {
	for _, v := range args {
		if strings.HasPrefix(v, "--reject=") {
			return v[9:]
		}
		if strings.HasPrefix(v, "-R=") {
			return v[3:]
		}
	}

	return ""
}

func get_Exclude(args []string) string {
	for _, v := range args {
		if strings.HasPrefix(v, "--exclude=") {
			return v[9:]
		}
		if strings.HasPrefix(v, "-X=") {
			return v[3:]
		}
	}

	return ""
}

func getBackground(args []string) string {
	for _, v := range args {
		if strings.HasPrefix(v, "-b=") {
			return v[3:]
		}
	}
	return ""
}
