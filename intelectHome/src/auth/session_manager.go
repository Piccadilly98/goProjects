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

type sessionManager struct {
	sessionIDcounter atomic.Int64
	sessionByJWT     map[string]*models.JWTinfo
	blackListJwtID   map[string]bool
	mtx              sync.Mutex
}

func MakeSessionManager() *sessionManager {
	sm := &sessionManager{
		sessionByJWT:   make(map[string]*models.JWTinfo),
		blackListJwtID: make(map[string]bool),
	}
	return sm
}

func (s *sessionManager) checkActiveSession(hash string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	info, ok := s.sessionByJWT[hash]
	if !ok {
		return false
	}

	if time.Now().After(info.Exp) {
		s.mtx.Lock()
		defer s.mtx.Unlock()
		delete(s.sessionByJWT, hash)
		return false
	}

	return true
}

func (s *sessionManager) CheckActiveSessionLogin(login string) (bool, string) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	for k, v := range s.sessionByJWT {
		if v.Login == login {
			return true, k
		}
	}
	return false, ""
}

func (s *sessionManager) CheckBlackListJWT(jwtHash string) bool {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	_, ok := s.blackListJwtID[jwtHash]
	return ok
}

func (s *sessionManager) NewSession(login string, role string, token string, exp time.Duration, id int64) bool {
	hash := s.hashToken(token)
	if s.checkActiveSession(hash) {
		return false
	}
	if s.CheckBlackListJWT(hash) {
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
	for k, v := range s.sessionByJWT {
		fmt.Println(k, v)
	}
	s.mtx.Unlock()
	return true
}

func (s *sessionManager) hashToken(token string) string {
	hash := sha256.New()
	hash.Write([]byte(token))
	return hex.EncodeToString(hash.Sum(nil))
}
