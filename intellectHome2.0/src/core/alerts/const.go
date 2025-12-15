package alerts

import "github.com/Piccadilly98/goProjects/intellectHome2.0/src/core/events"

const (
	TopicforBoardInfoChecker   = events.TopicBoardInfoUpdateDTO
	NameForBoardInfoChecker    = "board_info_checker"
	TopicForBoardStatusChecker = events.TopicBoardsStatusUpdate
	NameForBoardStatusChecker  = "board_status_checker"
)
