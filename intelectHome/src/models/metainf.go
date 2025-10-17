package models

type Metainf struct {
	UserInfo
	TokenID int64
	Exp     int64
}
