package main

import (
	"os"
	"wget/pkg/downloader"
	"wget/pkg/presenter"
)

// import "flag"


func main() {
	// 1. parse flags TODO write parser
	url := os.Args[1]

	// 2. init downloader
	presenter := presenter.NewCLIPresenter()

	downloader := downloader.NewDownloader(url, presenter)
	downloader.Download()
	//
}
