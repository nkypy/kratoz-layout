package auth

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/spf13/cast"

	authv1 "github.com/go-kratos/kratos-layout/api/auth/service/v1"
)

type VerifyService interface {
	VerifyToken(ctx context.Context, req *authv1.VerifyTokenReq) (*authv1.VerifyTokenResp, error)
}

const (
	tokenHeader          = "Authorization"    // Token header
	newTokenHeader       = "New-Token"        // 新 Token header
	currentUserSidHeader = "Current-User-Sid" // 当前用户 Sid header
)

var currentUserKey struct{}

type CurrentUser struct {
	Sid uint64
}

func JWTAuth(s VerifyService) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				result, err := s.VerifyToken(ctx, &authv1.VerifyTokenReq{
					Token:          tr.RequestHeader().Get(tokenHeader),
					CurrentUserSid: cast.ToUint64(tr.RequestHeader().Get(currentUserSidHeader)),
				})
				if err != nil {
					return nil, err
				}

				if result.NewToken != "" {
					tr.ReplyHeader().Set(newTokenHeader, result.NewToken)
				}

				// put CurrentUser into ctx
				ctx = WithContext(ctx, &CurrentUser{Sid: result.UserSid})
			}
			return handler(ctx, req)
		}
	}
}

func FromContext(ctx context.Context) *CurrentUser {
	return ctx.Value(currentUserKey).(*CurrentUser)
}

func WithContext(ctx context.Context, user *CurrentUser) context.Context {
	return context.WithValue(ctx, currentUserKey, user)
}
