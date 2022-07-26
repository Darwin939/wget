package presenter


type presenter interface {
	ShowStartTime()
	ShowRequestStatus()
	ShowContentSize()
	ShowName()

	ShowFinishTime()
}


type CLIPresenter struct {

}


func NewCLIPresenter() presenter {

	return &CLIPresenter{}
}
