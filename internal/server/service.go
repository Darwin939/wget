package server

import (
	"context"
	"os"
	"wget/internal/service"
	"wget/internal/service/downloader"
	"wget/internal/service/presenter"
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
	url := os.Args[1]

	// 2. init downloader
	presenter := presenter.NewCLIPresenter()

	downloader := downloader.NewDownloader(url, presenter)
	downloader.Download()
}
