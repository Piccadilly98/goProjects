package models

import "time"

type JWTinfo struct {
	Login string
	Role  string
	Exp   time.Time
	JwtID int64
}
