package service

import (
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
)

type Downloader interface {
	Download()
}

type Presenter interface {
	ShowStartTime()
	ShowRequestStatus(statusCode int)
	ShowContentSize(contentSize int64)
	ShowName(fullFilePath string)
	ShowProgress()
	ShowFinishTime([]string)
	GetBar(ContentLength int64) *progressbar.ProgressBar
}

type Mirrorer interface {
	CreateMirror()
}

type HTTPClient interface {
	SendHttp1(method string, url string, body io.Reader) (*http.Response, error)
}
