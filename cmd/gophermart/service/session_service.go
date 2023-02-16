package service

import (
	"github.com/fev0ks/ydx-goadv-tpl/model"
	"github.com/fev0ks/ydx-goadv-tpl/storage"
	"github.com/google/uuid"
	"time"
)

type SessionService interface {
	GetSession(sessionToken string) *model.Session
	CreateSession(userId int) (string, time.Time)
	DeleteSession(username string)
}

type sessionService struct {
	sessionStorage  storage.SessionStorage
	sessionLifetime time.Duration
}

func NewSessionService(sessionStorage storage.SessionStorage, sessionLifetime time.Duration) SessionService {
	return &sessionService{sessionStorage, sessionLifetime}
}

func (s sessionService) CreateSession(userId int) (string, time.Time) {
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(s.sessionLifetime)
	s.sessionStorage.SaveSession(sessionToken, &model.Session{
		UserId: userId,
		Expiry: expiresAt,
	})
	return sessionToken, expiresAt
}

func (s sessionService) GetSession(sessionToken string) *model.Session {
	return s.sessionStorage.GetSession(sessionToken)
}

func (s sessionService) DeleteSession(sessionToken string) {
	s.sessionStorage.DeleteSession(sessionToken)
}
