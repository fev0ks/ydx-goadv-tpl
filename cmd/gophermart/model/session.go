package model

import "time"

type Session struct {
	UserId int
	Expiry time.Time
}

func (s *Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
