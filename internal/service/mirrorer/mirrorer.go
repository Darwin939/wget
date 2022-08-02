package mirrorer

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
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

	if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
		return err
	}
	if err := m.parse(url, filePath, "index.html"); err != nil {
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
	for _, localPath := range localPaths {
		localPath := strings.TrimPrefix(localPath, url)
		fmt.Println("65url: ", url+"/"+localPath)
		err = m.download(url+"/"+localPath, filePath, localPath)
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
