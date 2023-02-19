package storage

import (
	"github.com/fev0ks/ydx-goadv-tpl/model"
)

// SessionStorage TODO use redis
type SessionStorage interface {
	GetSession(sessionToken string) *model.Session
	DeleteSession(sessionToken string)
	SaveSession(sessionToken string, session *model.Session)
}

type sessionStorage struct {
	storage map[string]*model.Session
}

func NewSessionStorage() SessionStorage {
	return &sessionStorage{storage: make(map[string]*model.Session, 0)}
}

func (s sessionStorage) GetSession(sessionToken string) *model.Session {
	return s.storage[sessionToken]
}

func (s sessionStorage) DeleteSession(sessionToken string) {
	delete(s.storage, sessionToken)
}

func (s sessionStorage) SaveSession(sessionToken string, session *model.Session) {
	s.storage[sessionToken] = session
}
