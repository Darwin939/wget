package config

type MirrorConf struct {
	Selectors map[string]string
}

type Config struct {
	Mirror MirrorConf
}

func GetConfig() Config {
	return Config{
		Mirror: MirrorConf{
			Selectors: map[string]string{
				"link":   "href",
				"img":    "src",
				"a":      "href",
				"script": "src",
			},
		},
	}
}
