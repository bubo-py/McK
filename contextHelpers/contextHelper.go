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

func RetrieveLoginFromContext(ctx context.Context) string {
	login := ctx.Value(loginKey).(string)
	return login
}

func RetrieveTimezoneFromContext(ctx context.Context) string {
	timezone := ctx.Value(timezoneKey).(string)
	return timezone
}
