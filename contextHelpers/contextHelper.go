package contextHelpers

import "context"

type contextKey string

var (
	loginKey    = contextKey("login")
	timezoneKey = contextKey("timezone")
)

func WriteLoginToContext(ctx context.Context, value string) context.Context {
	ctxWithData := context.WithValue(ctx, loginKey, value)
	return ctxWithData
}

func WriteTimezoneToContext(ctx context.Context, value string) context.Context {
	ctxWithData := context.WithValue(ctx, timezoneKey, value)
	return ctxWithData
}

func RetrieveLoginFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(loginKey).(string)
	return login, ok
}

func RetrieveTimezoneFromContext(ctx context.Context) (string, bool) {
	timezone, ok := ctx.Value(timezoneKey).(string)
	return timezone, ok
}
