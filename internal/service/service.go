package service

import "github.com/schollz/progressbar/v3"

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
