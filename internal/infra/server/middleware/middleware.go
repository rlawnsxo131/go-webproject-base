package middleware

const _middlewarePrefix = "internal/infra/server/middleware"

type contextKey struct {
	prefix string
	name   string
}

func newContextKey(name string) *contextKey {
	return &contextKey{_middlewarePrefix, name}
}

func (k *contextKey) String() string {
	return k.prefix + ":" + k.name
}
