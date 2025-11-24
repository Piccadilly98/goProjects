package dto

type UpdateBoardDto struct {
	Name  *string `json:"board_name"`
	Type  *string `json:"board_type"`
	State *string `json:"board_state"`
}

func (u *UpdateBoardDto) Validate() bool {
	return u.Name != nil || u.State != nil || u.Type != nil
}
