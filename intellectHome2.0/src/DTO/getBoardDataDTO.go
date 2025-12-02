package dto

import (
	"time"
)

type GetBoardDataDto struct {
	BoardId     *string    `json:"board_id"`
	Name        *string    `json:"board_name"`
	BoardType   *string    `json:"board_type"`
	BoardState  *string    `json:"board_state"`
	CreatedDate *time.Time `json:"created_date"`
	UpdatedDate *time.Time `json:"updated_date"`
}
