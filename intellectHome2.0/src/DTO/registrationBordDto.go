package dto

type RegistrationBoardDTO struct {
	BoardId    string  `json:"board_id"`
	Name       *string `json:"board_name"`
	BoardType  *string `json:"board_type"`
	BoardState *string `json:"board_state"`
}

func (u *RegistrationBoardDTO) ValidateAndDefault() bool {
	if u.BoardId == "" {
		return false
	}
	if u.BoardState == nil || *u.BoardState == "" {
		str := "registred"
		u.BoardState = &str
	}
	if u.BoardState != nil && (*u.BoardState != "registred" && *u.BoardState != "lost" &&
		*u.BoardState != "offline" && *u.BoardState != "active") {
		return false
	}
	if u.BoardType == nil || *u.BoardType == "" {
		str := "esp32_all_task"
		u.BoardType = &str
	}
	return true
}
