package consts

var UserCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}
