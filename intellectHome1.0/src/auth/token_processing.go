package auth

import (
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/models"
	"github.com/Piccadilly98/goProjects/intellectHome1.0/src/storage"
	"github.com/golang-jwt/jwt/v5"
)

type TokenWorker struct {
	tokenIDCount atomic.Int64
}

func (t *TokenWorker) CreateToken(login, role string, exp time.Duration) (string, int64, error) {
	t.tokenIDCount.Add(1)
	claims := &models.ClaimsJSON{
		Role:    role,
		TokenID: t.tokenIDCount.Load(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   login,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", 0, err
	}
	return tokenStr, claims.TokenID, nil
}

func ValidateToken(tokenStr string, stor *storage.Storage) (bool, *models.ClaimsJSON) {
	claims := &models.ClaimsJSON{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Println(err)
		return false, nil
	}
	if !token.Valid {
		return false, nil
	}
	roles := stor.GetAllRoles()
	if roles == nil {
		return false, nil
	}
	for _, v := range roles {
		if claims.Role == v {
			return true, claims
		}
	}
	return false, nil
}

func (t *TokenWorker) TokenToJSON(tokenStr string) ([]byte, error) {
	resp := &models.TokenResponseJSON{}
	claims := &models.ClaimsJSON{}
	resp.AccessToken = tokenStr
	resp.TokenType = "Bearer"
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid parsing")
	}
	resp.Role = claims.Role
	resp.TokenID = int(claims.TokenID)
	resp.UserID = claims.Subject
	resp.ExpirisIn = time.Until(claims.ExpiresAt.Time) / time.Second
	b, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}
