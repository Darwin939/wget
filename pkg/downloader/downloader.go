package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

type DownloaderI interface {
	Download()
	get()
	save()
	log()
}

type Downloader struct {
	Url          string
	IsFilename   bool
	IsPathPassed bool
	IsSpeedLimit bool
	Filename     string
	Path         string
	SpeedLimit   string
}

func NewDownloader(url string) DownloaderI {

	return &Downloader{Url: url}
}

func (d *Downloader) Download() {
	// if is file
	if !d.IsFilename {
		_, filename := path.Split(d.Url)
		d.Filename = filename
	}

	file, err := os.Create(d.Filename)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	resp, err := http.Get(d.Url)
	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Println(err)
	}


	if err != nil {
		log.Println(err)
	}
}

func (d *Downloader) get() {

}

func (d *Downloader) save() {}

func (d *Downloader) log() {}
