package authmw

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"bitbucket.org/efishery/go-efishery/helper"
	"bitbucket.org/efishery/go-efishery/libs/token"
	libtoken "bitbucket.org/efishery/go-efishery/libs/token"
	"github.com/dgrijalva/jwt-go"
)

// HTTPAuthMiddleware is middleware to validate token is valid from server
func HTTPAuthMiddleware(next http.Handler, signatureKey string) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var tokenString string
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.Contains(authHeader, bearerSchema) {
			tokenString = authHeader[len(bearerSchema):]
		}

		if tokenString == "" {
			c, _ := r.Cookie("token")
			if c != nil {
				tokenString = c.Value
			}
		}

		if tokenString == "" {
			tokenString = r.URL.Query().Get("token")
		}

		tokenString = strings.TrimSpace(tokenString)
		token, err := token.NewService(signatureKey).Validate(tokenString)
		if err != nil {
			helper.HTTPWriteResponse(rw, err, http.StatusUnauthorized)
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			// currently using marshal unmarshal
			// readable but need tunning performace
			byteClaims, err := json.Marshal(claims)
			if err != nil {
				helper.HTTPWriteResponse(rw, err, http.StatusUnauthorized)
				return
			}
			myClaim := libtoken.AuthServiceJWT{}
			err = json.Unmarshal(byteClaims, &myClaim)
			if err != nil {
				helper.HTTPWriteResponse(rw, err, http.StatusUnauthorized)
				return
			}

			// inject context token
			ctx := context.WithValue(r.Context(), contextClaimKey, myClaim)
			next.ServeHTTP(rw, r.WithContext(ctx))
		} else {
			helper.HTTPWriteResponse(rw, err, http.StatusUnauthorized)
		}

	})
}

// GetClaim is get token claim from user
func GetClaim(r *http.Request) token.AuthServiceJWT {
	v := r.Context().Value(contextClaimKey)
	if v == nil {
		return libtoken.AuthServiceJWT{}
	}

	out, ok := v.(libtoken.AuthServiceJWT)
	if !ok {
		return libtoken.AuthServiceJWT{}
	}

	return out
}
