package downloader

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"wget/internal/service"

	"github.com/mxk/go-flowrate/flowrate"
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

	body := resp.Body

	if d.IsSpeedLimit {
		speedLimit, err := calculateSpeedLimit(d.SpeedLimit)
		
		if err != nil {
			log.Println(err)
		}
		body = flowrate.NewReader(resp.Body, speedLimit)
	}

	_, err = io.Copy(io.MultiWriter(file, d.Presenter.GetBar(resp.ContentLength)), body)
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


// calculates bytes per second
func calculateSpeedLimit(speed string) (int64, error) {
	if strings.HasPrefix(speed, "k") || strings.HasPrefix(speed, "K") {
		num, err := strconv.Atoi(speed[:len(speed)-1])
		if err != nil {
			return -1 , errors.New("wrong speed limit argument")
		}

		return int64(num * 1024), nil
	} else if strings.HasPrefix(speed, "m") || strings.HasPrefix(speed, "M") {
		num, err := strconv.Atoi(speed[:len(speed)-1])
		if err != nil {
			return -1 , errors.New("wrong speed limit argument")
		}

		return int64(num * 1024 * 1024), nil
	}

	return -1, errors.New("wrong speed limit argument")

}
