package auth

import (
	"encoding/json"
	"io"
	"os"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

func ValidateLoginData(body io.ReadCloser, stor *storage.Storage) (bool, string, string) {
	b, err := io.ReadAll(body)
	if err != nil {
		return false, "", ""
	}
	login := &models.LoginJSON{}
	roles := stor.GetAllRoles()
	if roles == nil {
		return false, "", ""
	}
	err = json.Unmarshal(b, login)
	if err != nil {
		return false, "", ""
	}
	for _, v := range roles {
		env := os.Getenv(v + "_LOGIN")
		if login.Login == env && login.Password == os.Getenv(v+"_PASSWORD") {
			login.Role = v
			return true, login.Login, login.Role
		}
	}
	return false, "", ""
}
