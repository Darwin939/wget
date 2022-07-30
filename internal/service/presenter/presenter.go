package presenter

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
)

type CLIPresenter struct {
}

func NewCLIPresenter() *CLIPresenter {

	return &CLIPresenter{}
}

func (c *CLIPresenter) ShowStartTime() {
	currentTime := time.Now()

	fmt.Printf("start at %v\n", currentTime.Format("2006-01-02 15:04:05"))
}

func (c *CLIPresenter) ShowRequestStatus(statusCode int) {
	if statusCode >= 200 && statusCode <= 299 {
		fmt.Printf("sending request, awaiting response... status %v OK\n", statusCode)
	} else {
		fmt.Printf("sending request, awaiting response... status %v something goes wrong\n", statusCode)
	}

}

func (c *CLIPresenter) ShowContentSize(contentSize int64) {
	mb := (float64(contentSize) / 1024) / 1024
	fmt.Printf("content size:  %v [±%fMB]\n", contentSize, mb)

}

func (c *CLIPresenter) ShowName(fullFilePath string) {

	fmt.Printf("saving file to: %v\n", fullFilePath)
}

func (c *CLIPresenter) ShowFinishTime(files []string) {
	currentTime := time.Now()

	fmt.Printf("Downloaded %v\n", files)

	fmt.Printf("finished at %v\n", currentTime.Format("2006-01-02 15:04:05"))

}

func (c *CLIPresenter) ShowProgress() {

}

func (c *CLIPresenter) GetBar(ContentLength int64) *progressbar.ProgressBar {
	return progressbar.DefaultBytes(
		ContentLength,
		"downloading",
	)
}
