package utils

import "context"

var (
	tokenContextKey = "conn-context"
)

type ConnContext struct {
	Token     string
	Platform  string
	SessionId string
}

func SetConnContext(ctx context.Context, token string, sessionId string, platform string) context.Context {
	return context.WithValue(
		ctx, tokenContextKey, &ConnContext{
			Token:     token,
			Platform:  platform,
			SessionId: sessionId,
		},
	)
}

func MustGetConnContext(ctx context.Context) *ConnContext {
	connContext, ok := ctx.Value(tokenContextKey).(*ConnContext)
	if !ok {
		panic("ConnContext not found")
	}

	return connContext
}
