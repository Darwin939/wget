package downloader

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"wget/internal/service"
)

type Downloader struct {
	Url          string
	IsFilename   bool
	IsPathPassed bool
	IsSpeedLimit bool
	Filename     string
	Path         string
	SpeedLimit   string
	fullFilePath string

	Presenter service.Presenter
}

func NewDownloader(url string, presenter service.Presenter) *Downloader {

	return &Downloader{Url: url, Presenter: presenter}
}

func (d *Downloader) Download() {
	// if is file
	d.Presenter.ShowStartTime()

	if !d.IsFilename {
		_, filename := path.Split(d.Url)
		d.Filename = filename
	}

	file, err := os.Create(d.Filename)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	resp := d.get()

	d.Presenter.ShowRequestStatus(resp.StatusCode)
	d.Presenter.ShowContentSize(resp.ContentLength)

	defer resp.Body.Close()

	_, err = io.Copy(io.MultiWriter(file, d.Presenter.GetBar(resp.ContentLength)), resp.Body)
	if err != nil {
		log.Println(err)
	}
	d.Presenter.ShowName(d.generateFileFullPath())

	if err != nil {
		log.Println(err)
	}
	d.Presenter.ShowFinishTime([]string{d.Url})
}

func (d *Downloader) get() *http.Response {
	resp, err := http.Get(d.Url)
	if err != nil {
		log.Println(err)
	}
	return resp
}

func (d *Downloader) save() {}

func (d *Downloader) log() {}

func (d *Downloader) generateFileFullPath() string {

	return ""
}
