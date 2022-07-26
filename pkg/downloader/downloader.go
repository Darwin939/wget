package downloader



type DownloaderI interface {
	Download()
	get()
	save()
	log()
}


type Downloader struct {
    Url string
    IsOwnName bool
    IsPathPassed bool
    IsSpeedLimit bool
    OwnName string
    Path string
    SpeedLimit string
}

func NewDownloader(url string) DownloaderI{

	return &Downloader{Url: url}
}


func (d *Downloader) Download() {

}

func (d *Downloader) get() {

}

func (d *Downloader) save() {}

func (d *Downloader) log() {}
