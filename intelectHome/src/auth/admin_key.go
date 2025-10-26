package auth

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Piccadilly98/goProjects/intelectHome/src/storage"
)

func PrintAdminKey(url string, tw *TokenWorker, sm *SessionManager, stor *storage.Storage) string {
	token, id, err := tw.CreateToken(os.Getenv("ADMIN_LOGIN"), "ADMIN", 100*time.Hour)
	urlResult := ""
	if err != nil {
		log.Fatal(err)
	}
	if !sm.NewSession(os.Getenv("ADMIN_LOGIN"), "ADMIN", token, 100*time.Hour, id) {
		log.Fatal("error new session")
	}
	hash := sm.getValidHashToken(token)
	if hash == "" {
		log.Fatal("no session")
	}
	fmt.Printf("Admin token: %s\nHash: %s\n", token, hash)
	if strings.HasPrefix(url, ":") {
		url = "localhost" + url
	}
	urlResult = "http://" + url + "/quick-auth-admin?hash=" + hash
	fmt.Printf("\n\nAdd cookie in:\n%s\n", urlResult)
	return token
}
