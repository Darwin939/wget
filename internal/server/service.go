package server

import (
	"context"
	"wget/internal/service"
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

}
