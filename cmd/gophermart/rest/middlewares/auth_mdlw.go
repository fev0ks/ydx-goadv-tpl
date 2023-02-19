package middlewares

import (
	"context"
	"github.com/fev0ks/ydx-goadv-tpl/model/consts"
	"github.com/fev0ks/ydx-goadv-tpl/service"
	"net/http"
)

type SessionTokenValidator interface {
	ValidateSessionToken(next http.Handler) http.Handler
}

type authMiddleware struct {
	SessionService service.SessionService
}

func NewAuthMiddleware(sessionService service.SessionService) SessionTokenValidator {
	return &authMiddleware{
		SessionService: sessionService,
	}
}

func (am *authMiddleware) ValidateSessionToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {

				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := cookie.Value
		session := am.SessionService.GetSession(sessionToken)
		if session == nil || session.IsExpired() {
			am.SessionService.DeleteSession(sessionToken)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), consts.UserIDCtxKey, session.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
