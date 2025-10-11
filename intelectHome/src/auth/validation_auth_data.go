package auth

import (
	"encoding/json"
	"io"
	"os"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
)

func ValidateLoginData(body io.ReadCloser) (bool, string, string) {
	b, err := io.ReadAll(body)
	if err != nil {
		return false, "", ""
	}
	names := []string{"ESP32_1", "ADMIN"}
	login := &models.LoginJSON{}
	err = json.Unmarshal(b, login)
	if err != nil {
		return false, "", ""
	}
	for _, v := range names {
		env := os.Getenv(v + "_LOGIN")
		if login.Login == env && login.Password == os.Getenv(v+"_PASSWORD") {
			login.Role = v
			return true, login.Login, login.Role
		}
	}
	return false, "", ""
}
