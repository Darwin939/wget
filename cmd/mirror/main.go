package main

import (
	"fmt"
	"os"
	"wget/internal/config"
	downloader2 "wget/internal/service/downloader"
	"wget/internal/service/mirrorer"
	"wget/internal/service/parser"
	"wget/internal/service/presenter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("arg should count >=2")
		return
	}
	flags := parser.NewFlags(os.Args)
	presenter := presenter.NewCLIPresenter()
	downloader := downloader2.NewDownloader(flags, presenter, 0)
	conf := config.GetConfig()
	service := mirrorer.NewMirrorer(conf.Mirror, downloader, os.Args[1], nil, nil)
	if err := service.CreateMirror(); err != nil {
		fmt.Println(err)
		return
	}
}
