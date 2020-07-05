package core

import "context"

const (
	MobContext = "mobCtx"
)

type APIContext struct {
	APIKey       string
	APIKeySecret string
	Email        string
	RequestID    string
}

func ContextWithValue(ctx context.Context, val interface{}) context.Context {
	return context.WithValue(ctx, MobContext, val)
}

func AddToContext(ctx context.Context, key, val interface{}) context.Context {
	return context.WithValue(ctx, key, val)
}

func getValue(ctx context.Context, key interface{}) interface{} {
	return ctx.Value(key)
}

func GetContextData(ctx context.Context) APIContext {
	return ctx.Value(MobContext).(APIContext)
}

func GetReuestID(ctx context.Context) string {
	return getValue(ctx, RequestID).(string)
}
