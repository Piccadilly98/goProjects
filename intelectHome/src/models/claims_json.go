package models

import "github.com/golang-jwt/jwt/v5"

type ClaimsJSON struct {
	Role    string `json:"role"`
	TokenID int64  `json:"tokenID"`
	jwt.RegisteredClaims
}
