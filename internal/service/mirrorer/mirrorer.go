package mirrorer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
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
	downloader    service.Downloader
	presenter     service.Presenter
}

func NewMirrorer(downloader service.Downloader, url string, excluded, reject []string) *Mirrorer {
	return &Mirrorer{
		excluded:   excluded,
		downloader: downloader,
		reject:     reject,
		url:        url,
	}
}

func (m *Mirrorer) CreateMirror() error {
	m.initRegex()
	url := validateURL(m.url)
	filePath := convertToPath(url)

	//if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
	//	return err
	//}
	if err := m.parse(url, filePath, "index.html"); err != nil {
		return err
	}
	return nil
}

func (m *Mirrorer) parse(url, filePath, name string) error {

	err := m.downloader.Download(url, filePath, name)
	if err != nil {
		return err
	}
	//err := m.download(url, filePath, name)
	//if err != nil {
	//	return err
	//}
	file, err := os.Open(path.Join(filePath, "index.html"))
	if err != nil {
		return err
	}
	//
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return err
	}
	selection := doc.Find("style")
	localPaths := FindPath(selection.Text())
	selectors := map[string]string{
		"link":   "href",
		"img":    "src",
		"a":      "href",
		"script": "src",
	}
	for selector, attr := range selectors {
		doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
			val, exists := selection.Attr(attr)
			if exists {
				//fmt.Println("href:", val)
				if isLocal, localURL := isLocalPath(val); isLocal {
					localPaths = append(localPaths, localURL)
				}
			}
		})
	}

	fmt.Println("localPaths:", localPaths)
	for _, localPath := range localPaths {
		localPath = strings.TrimPrefix(localPath, url)
		localFilePath := filePath
		fmt.Println("91url:", url, "filepath:", localFilePath)
		dir, filename := path.Split(localPath)

		localFilePath = path.Join(localFilePath, dir)
		err = m.downloader.Download(url+localPath, localFilePath, filename)
		if err != nil {
			return err
		}
	}

	//doc.Find("link").
	return nil
}

func (m *Mirrorer) download(url, filePath, name string) error {
	//m.presenter.ShowStartTime()
	//resp, err := m.cli.SendHttp1(http.MethodGet, url, nil)
	//if err != nil {
	//	return err
	//}
	//defer resp.Body.Close()
	//fmt.Println("filePath:", filePath, "name: ", name)
	//dir, _ := path.Split(filePath + name)
	//fmt.Println("dir:", dir)
	//if err = os.MkdirAll(dir, os.ModePerm); err != nil {
	//	return err
	//}
	//m.presenter.ShowRequestStatus(resp.StatusCode)
	//m.presenter.ShowContentSize(resp.ContentLength)
	//fmt.Println("91path:", path.Join(filePath, name))
	//file, err := os.Create(path.Join(filePath, name))
	//if err != nil {
	//	return err
	//}
	//defer file.Close()

	//_, err = io.Copy(io.MultiWriter(file, m.presenter.GetBar(resp.ContentLength)), resp.Body)
	//if err != nil {
	//	return err
	//}
	return nil
}

func (m *Mirrorer) initRegex() {
	//TODO find optimal way for creating regex
	m.excludedRegex = regexBuilder(`\.`, m.excluded)
	m.rejectRegex = regexBuilder(`\/`, m.reject)

}
