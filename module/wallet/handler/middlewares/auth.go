package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/herwando/mini-wallet/lib/common/commonerr"
	"github.com/herwando/mini-wallet/lib/common/writer"
	"github.com/herwando/mini-wallet/module/wallet/entity/model"
)

var (
	jwtKey                  = []byte("my_secret_key")
	writerWriteJSONAPIError = writer.WriteJSONAPIError
	CONTEXT_AUTH_DETAIL     = "AuthDetail"
)

type Module struct {
}

const authPrefix string = "Token "

func GetAuthDetailFromContext(ctx context.Context) (string, error) {
	v := ctx.Value(CONTEXT_AUTH_DETAIL)
	if v == nil {
		return "", errors.New("unable to get account from context")
	}
	return v.(string), nil
}

func (m *Module) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		tokenStr := getBearerToken(r, "Authorization")
		if tokenStr == "" {
			writerWriteJSONAPIError(ctx, w, commonerr.SetNewBadRequest("Request invalid", "Header Authorization empty"))
			return
		}
		claims := &model.Claims{}

		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnauthorizedError("Unauthorized user", "You are not authorized to use this function"))
				return
			}
			writerWriteJSONAPIError(ctx, w, commonerr.SetDefaultNewBadRequest())
			return
		}

		if !tkn.Valid {
			writerWriteJSONAPIError(ctx, w, commonerr.SetNewUnauthorizedError("Unauthorized user", "You are not authorized to use this function"))
			return
		}

		ctx = context.WithValue(ctx, CONTEXT_AUTH_DETAIL, claims.CustomerXid)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getBearerToken(r *http.Request, key string) string {
	var token string
	if authHeader := r.Header.Get(key); authHeader != "" {
		token = strings.Replace(authHeader, authPrefix, "", 1)
	}
	return token
}
