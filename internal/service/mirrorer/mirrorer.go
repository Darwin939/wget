package mirrorer

import (
	"fmt"
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
	//if err := m.download(url, folders, filename); err != nil {
	//	return err
	//}
	if err := m.parse(url, folders, filename); err != nil {
		return err
	}
	return nil
}

func (m *Mirrorer) parse(url, filePath, name string) error {
	err := m.download(url, filePath, name)
	if err != nil {
		return err
	}
	file, err := os.Open(path.Join(filePath, "index.html"))
	if err != nil {
		return err
	}

	defer file.Close()
	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	//fmt.Println("b:", string(b))
	localPaths := FindPath(b)
	fmt.Println(localPaths)
	dir, _ := path.Split(url)
	for _, localPath := range localPaths {
		ldir, lfile := path.Split(localPath)
		err = m.download(dir+localPath, filePath+ldir, lfile)
		if err != nil {
			return err
		}
	}

	//doc.Find("link").
	return nil
}

func (m *Mirrorer) download(url, filePath, name string) error {
	m.presenter.ShowStartTime()
	resp, err := m.cli.SendHttp1(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	m.presenter.ShowRequestStatus(resp.StatusCode)
	m.presenter.ShowContentSize(resp.ContentLength)
	p := path.Join(filePath, name)
	stat, err := os.Stat(p)

	if err == nil && stat.IsDir() {
		name = "index.html"
	}
	file, err := os.Create(path.Join(filePath, name))
	if err != nil {
		return err
	}
	defer file.Close()

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
