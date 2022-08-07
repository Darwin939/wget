package mirrorer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"os"
	"path"
	"wget/internal/config"
	"wget/internal/service"
)

type Mirrorer struct {
	excludedRegex string
	rejectRegex   string
	url           string
	selectors     map[string]string
	excluded      []string
	reject        []string
	downloader    service.Downloader
	presenter     service.Presenter
}

func NewMirrorer(conf config.MirrorConf, downloader service.Downloader, url string, excluded, reject []string) *Mirrorer {
	return &Mirrorer{
		excluded:   excluded,
		downloader: downloader,
		selectors:  conf.Selectors,
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

	if err := m.downloader.Download(url, filePath, name); err != nil {
		return err
	}

	file, err := os.Open(path.Join(filePath, "index.html"))
	if err != nil {
		return err
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return err
	}
	localPaths := FindPath(doc.Find("style").Text())

	for selector, attr := range m.selectors {
		doc.Find(selector).Each(func(i int, selection *goquery.Selection) {
			val, exists := selection.Attr(attr)
			if exists {
				if isLocal, localURL := isLocalPath(val); isLocal {
					localPaths = append(localPaths, localURL)
				}
			}
		})
	}

	fmt.Println("localPaths:", localPaths)
	for _, localPath := range localPaths {
		dir, filename := path.Split(localPath)
		localFilePath := filePath
		localFilePath = path.Join(localFilePath, dir)
		err = m.downloader.Download(url+localPath, localFilePath, filename)
		if err != nil {
			return err
		}
	}

	//doc.Find("link").
	return nil
}

func (m *Mirrorer) initRegex() {
	//TODO find optimal way for creating regex
	m.excludedRegex = regexBuilder(`\.`, m.excluded)
	m.rejectRegex = regexBuilder(`\/`, m.reject)

}
