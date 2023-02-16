package consts

var UserIdCtxKey = &contextKey{"userId"}

type contextKey struct {
	name string
}
