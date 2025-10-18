package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Piccadilly98/goProjects/intelectHome/src/models"
)

type SessionManager struct {
	sessionIDcounter atomic.Int64
	sessionByJWT     map[string]*models.JWTinfo
	blackListJwtID   map[string]bool
	mtx              sync.Mutex
}

func MakeSessionManager() *SessionManager {
	sm := &SessionManager{
		sessionByJWT:   make(map[string]*models.JWTinfo),
		blackListJwtID: make(map[string]bool),
	}
	sm.blackListJwtID["fea7353b9c1b5449bb59071c62f908d2b5d5ce5ad654bdbfde31dc2ff0e4a4b6"] = true
	return sm
}

func (s *SessionManager) checkActiveSession(hash string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	info, ok := s.sessionByJWT[hash]
	if !ok {
		return false
	}

	if time.Now().After(info.Exp) {
		delete(s.sessionByJWT, hash)
		return false
	}

	return true
}

func (s *SessionManager) checkActiveSessionLogin(login string) (bool, string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for k, v := range s.sessionByJWT {
		if v.Login == login {
			return true, k
		}
	}
	return false, ""
}

func (s *SessionManager) checkBlackListJWT(jwtHash string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.blackListJwtID[jwtHash]
	return ok
}

func (s *SessionManager) NewSession(login string, role string, token string, exp time.Duration, id int64) bool {
	hash := s.hashToken(token)
	if s.checkActiveSession(hash) {
		return true
	}
	if ok, _ := s.checkActiveSessionLogin(login); ok {
		return false
	}
	if s.checkBlackListJWT(hash) {
		return false
	}
	s.mtx.Lock()
	s.sessionIDcounter.Add(1)
	s.sessionByJWT[hash] = &models.JWTinfo{
		JwtID: id,
		Login: login,
		Role:  role,
		Exp:   time.Now().Add(exp),
	}
	s.mtx.Unlock()
	return true
}

func (s *SessionManager) hashToken(token string) string {
	hash := sha256.New()
	hash.Write([]byte(token))
	return hex.EncodeToString(hash.Sum(nil))
}

func (s *SessionManager) CheckTokenValid(token string, claims *models.ClaimsJSON) (bool, error) {
	var err error
	hash := s.hashToken(token)
	if s.checkBlackListJWT(hash) {
		err = fmt.Errorf("jwt in BL, hash: %s", hash)
		return false, err
	}
	if !s.checkActiveSession(hash) {
		err = fmt.Errorf("jwt not have active session, need /login")
		return false, err
	}
	if ok, _ := s.checkActiveSessionLogin(claims.Subject); !ok {
		err = fmt.Errorf("login not have active session, need /login")
		return false, err
	}
	return true, nil
}
