package consts

var UserIdCtxKey = &contextKey{"userID"}

type contextKey struct {
	name string
}
