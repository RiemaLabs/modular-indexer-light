package httputl

import (
	"context"

	"github.com/google/uuid"
)

var reqIDKey struct{}

func TODO() context.Context {
	return WithRequestID(context.TODO(), "")
}

func WithRequestID(ctx context.Context, reqID string) context.Context {
	if reqID == "" {
		reqID = uuid.NewString()
	}
	return context.WithValue(ctx, reqIDKey, reqID)
}

func RequestID(ctx context.Context) string {
	if reqID, ok := ctx.Value(reqIDKey).(string); ok {
		return reqID
	}
	return ""
}
