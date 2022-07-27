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
	"wget/internal/service/parser"

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
	IsSaveFrom   bool
	IsMirror     bool
	IsReject     bool
	IsExclude    bool
	SaveFrom     string
	Reject       string
	Exclude      string

	Presenter service.Presenter
}

func NewDownloader(flag *parser.Flags, presenter service.Presenter) *Downloader {

	return &Downloader{Url: flag.Url,
		Presenter:    presenter,
		IsFilename:   flag.IsFilename,
		IsPathPassed: flag.IsPathPassed,
		IsSpeedLimit: flag.IsSpeedLimit,
		Filename:     flag.Filename,
		Path:         flag.Path,
		IsMirror:     flag.IsMirror,
		IsSaveFrom:   flag.IsSaveFrom,
		IsReject:     flag.IsReject,
		IsExclude:    flag.IsExclude,
		SaveFrom:     flag.SaveFrom,
		Reject:       flag.Reject,
		Exclude:      flag.Exclude,
		SpeedLimit:   flag.SpeedLimit,
		// TODO а есть удобный метод распаковки
	}
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
	if strings.HasSuffix(speed, "k") || strings.HasSuffix(speed, "K") {
		num, err := strconv.Atoi(speed[:len(speed)-1])
		if err != nil {
			return -1, errors.New("wrong speed limit argument")
		}

		return int64(num * 1024), nil
	} else if strings.HasSuffix(speed, "m") || strings.HasSuffix(speed, "M") {
		num, err := strconv.Atoi(speed[:len(speed)-1])
		if err != nil {
			return -1, errors.New("wrong speed limit argument")
		}

		return int64(num * 1024 * 1024), nil
	}

	return -1, errors.New("wrong speed limit argument")

}
