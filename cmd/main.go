package main

import "wget/pkg/downloader"
// import "flag"
import "os"


func main() {
	// 1. parse flags TODO write parser
	url := os.Args[1]

	// 2. init downloader
	Downloader := downloader.NewDownloader(url)
	Downloader.Download()
	//
}
