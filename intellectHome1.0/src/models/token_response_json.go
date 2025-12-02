package models

import "time"

type TokenResponseJSON struct {
	AccessToken string        `json:"accessToken"`
	TokenType   string        `json:"tokenType"`
	ExpirisIn   time.Duration `json:"exprisIn"`
	TokenID     int           `json:"tokenID"`
	UserInfo
}

type UserInfo struct {
	UserID string `json:"userID"`
	Role   string `json:"role"`
}
