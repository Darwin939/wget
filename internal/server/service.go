package server

import (
	"context"
	"fmt"
	"os"
	"wget/internal/service"
	"wget/internal/service/downloader"
	"wget/internal/service/presenter"

	"wget/internal/service/parser"
)

type App struct {
	downloader service.Downloader
	presenter  service.Presenter
}

func NewApp() *App {
	return &App{}
}

func (a *App) Initialize() {

}

func (a *App) Run(ctx context.Context) {

	flags := parser.NewFlags(os.Args)

	// 2. init downloader
	presenter := presenter.NewCLIPresenter()

	downloader := downloader.NewDownloader(flags, presenter)
	switch {
	case downloader.IsSaveFrom: // -i flag
		fmt.Println("case 1")
		downloader.DownloadFromFile()
	case downloader.IsMirror:
		fmt.Println("case 2")

	default:
		fmt.Println("case 3")
		downloader.Download(flags.Url)
	}
}
