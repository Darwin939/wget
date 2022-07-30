package mirrorer

import (
	"fmt"
	"github.com/gabriel-vasile/mimetype"
	"io"
	"net/http"
	"os"
	"path"
	"wget/internal/service"
)

type Mirrorer struct {
	excludedRegex string
	rejectRegex   string
	url           string
	excluded      []string
	reject        []string
	cli           service.HTTPClient
	presenter     service.Presenter
}

func NewMirrorer(cli service.HTTPClient, presenter service.Presenter, url string, excluded, reject []string) *Mirrorer {
	return &Mirrorer{
		excluded:  excluded,
		reject:    reject,
		url:       url,
		cli:       cli,
		presenter: presenter,
	}
}

func (m *Mirrorer) CreateMirror() error {
	m.initRegex()
	url := validateURL(m.url)
	filePath := convertToPath(url)
	dir, filename := path.Split(filePath)
	//fmt.Println("dir:", dir, "filename:", filename)
	var folders = dir
	if dir == "" {
		folders = filename
	}
	if err := os.MkdirAll(folders, os.ModePerm); err != nil {
		return err
	}
	if err := m.download(url, folders, filename); err != nil {
		return err
	}
	return nil
}
func (m *Mirrorer) download(url, path, name string) error {
	m.presenter.ShowStartTime()
	resp, err := m.cli.SendHttp1(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	m.presenter.ShowRequestStatus(resp.StatusCode)
	m.presenter.ShowContentSize(resp.ContentLength)
	mtype, err := mimetype.DetectReader(resp.Body)
	if err != nil {
		return err
	}
	if mtype.Is("text/html") {
		name = "index.html"
	}

	file, err := os.Create(fmt.Sprintf("%s/%s", path, name))
	if err != nil {
		return err
	}

	_, err = io.Copy(io.MultiWriter(file, m.presenter.GetBar(resp.ContentLength)), resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (m *Mirrorer) initRegex() {
	//TODO find optimal way for creating regex
	m.excludedRegex = regexBuilder(`\.`, m.excluded)
	m.rejectRegex = regexBuilder(`\/`, m.reject)

}
