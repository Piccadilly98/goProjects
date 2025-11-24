package dto

type UploadBoardDataDto struct {
	BoardId    *string `json:"board_id"`
	Name       *string `json:"board_name"`
	BoardType  *string `json:"board_type"`
	BoardState *string `json:"board_state"`
}

func (u *UploadBoardDataDto) ValidateAndDefault() bool {
	if u.BoardId == nil || *u.BoardId == "" {
		return false
	}
	if u.BoardState == nil || *u.BoardId == "" {
		str := "registred"
		u.BoardState = &str
	}
	if u.BoardType == nil || *u.BoardType == "" {
		str := "esp32_all_task"
		u.BoardType = &str
	}
	return true
}
