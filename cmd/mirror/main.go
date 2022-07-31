package main

import (
	"fmt"
	"os"
	"time"
	"wget/internal/service/client"
	"wget/internal/service/mirrorer"
	"wget/internal/service/presenter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("arg should count >=2")
		return
	}
	present := presenter.NewCLIPresenter()
	cli := client.NewClient(15 * time.Second)
	service := mirrorer.NewMirrorer(cli, present, os.Args[1], nil, nil)
	if err := service.CreateMirror(); err != nil {
		fmt.Println(err)
		return
	}
}
