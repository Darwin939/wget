package downloader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"wget/internal/service"
	"wget/internal/service/client"
	"wget/internal/service/parser"

	"github.com/gabriel-vasile/mimetype"
	"github.com/mxk/go-flowrate/flowrate"
)

type Downloader struct {
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
	IsBackground bool
	SaveFrom     string
	Reject       string
	Exclude      string

	Presenter service.Presenter

	cli service.HTTPClient
}

func NewDownloader(flag *parser.Flags, presenter service.Presenter, timeout time.Duration) *Downloader {

	return &Downloader{
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
		IsBackground: flag.IsBackground,

		cli: client.NewClient(timeout),
		// TODO а есть удобный метод распаковки
	}
}

func (d *Downloader) Download(Url, Path, FileName string) error {
	// if is file
	d.Presenter.ShowStartTime()

	resp, err := d.cli.SendHttp1(http.MethodGet, Url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	mtype, err := mimetype.DetectReader(resp.Body)
	if err != nil {
		return err
	}
	if mtype.Is("text/html") {
		FileName = "index.html"
	}

	if d.IsBackground {

	}

	if !d.IsFilename { // -O
		_, filename := path.Split(Url)
		d.Filename = filename
	}
	log.Printf("create file: %v\n", filepath.Join(Path, FileName))
	file, err := os.Create(filepath.Join(Path, FileName))
	if err != nil {
		return err
	}

	defer file.Close()

	resp = d.get(Url)

	d.Presenter.ShowRequestStatus(resp.StatusCode)
	d.Presenter.ShowContentSize(resp.ContentLength)

	defer resp.Body.Close()

	body := resp.Body

	if d.IsSpeedLimit {

		speedLimit, err := calculateSpeedLimit(d.SpeedLimit)

		if err != nil {
			return err
		}
		body = flowrate.NewReader(resp.Body, speedLimit)
	}

	_, err = io.Copy(io.MultiWriter(file, d.Presenter.GetBar(resp.ContentLength)), body)
	if err != nil {
		return err
	}
	d.Presenter.ShowName(d.generateFileFullPath())

	if err != nil {
		return err
	}
	d.Presenter.ShowFinishTime([]string{Url})
	return nil
}

func (d *Downloader) get(url string) *http.Response {

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	return resp
}

func (d *Downloader) save() {}

func (d *Downloader) log() {}

func (d *Downloader) generateFileFullPath() string {

	return filepath.Join(d.Path, d.Filename)
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

func (d *Downloader) DownloadFromFile() {
	var urls []string

	f, err := os.Open(d.SaveFrom)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	log.Printf("savefrom: %v\n", urls)

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	for i := 0; i < len(urls); i++ {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()
			// d.Url = urls[i]
			d.Download(url, d.Path, d.Filename)
		}(urls[i])
	}
	wg.Wait()
}
