package middlewares

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strings"

	keycloakclient "github.com/Slava02/ChatSupport/internal/clients/keycloak"
	"github.com/Slava02/ChatSupport/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/introspector_mock.gen.go -package=middlewaresmocks Introspector

const tokenCtxKey = "user-token"

var ErrNoRequiredResourceRole = errors.New("no required resource role")

type Introspector interface {
	IntrospectToken(ctx context.Context, token string) (*keycloakclient.IntrospectTokenResult, error)
}

// NewKeycloakTokenAuth returns a middleware that implements "active" authentication:
// each request is verified by the Keycloak server.
func NewKeycloakTokenAuth(introspector Introspector, resource, role string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Bearer",
		Validator: func(tokenStr string, eCtx echo.Context) (bool, error) {
			tokenStr = sanitize(tokenStr)

			token, err := introspector.IntrospectToken(eCtx.Request().Context(), tokenStr)
			if err != nil {
				return false, err
			}

			if !token.Active {
				return false, nil
			}

			var cl claims
			t, _, _ := new(jwt.Parser).ParseUnverified(tokenStr, &cl)
			if err != nil {
				return false, err
			}

			if err = cl.Valid(); err != nil {
				return false, err
			}

			if !cl.ResourceAccess.HasResourceRole(role, resource) {
				return false, echo.ErrForbidden.WithInternal(ErrNoRequiredResourceRole)
			}

			eCtx.Set(tokenCtxKey, t)

			return true, nil
		},
	})
}

func sanitize(t string) string {
	for _, ch := range []string{" ", ","} {
		t = strings.ReplaceAll(t, ch, "")
	}
	return t
}

func MustUserID(eCtx echo.Context) types.UserID {
	uid, ok := userID(eCtx)
	if !ok {
		panic("no user token in request context")
	}
	return uid
}

func userID(eCtx echo.Context) (types.UserID, bool) {
	t := eCtx.Get(tokenCtxKey)
	if t == nil {
		return types.UserIDNil, false
	}

	tt, ok := t.(*jwt.Token)
	if !ok {
		return types.UserIDNil, false
	}

	userIDProvider, ok := tt.Claims.(interface{ UserID() types.UserID })
	if !ok {
		return types.UserIDNil, false
	}
	return userIDProvider.UserID(), true
}
