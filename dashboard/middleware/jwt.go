package middleware

import (
	"net/http"
	"strings"

	"github.com/telkomdev/tob/dashboard/shared"
	"github.com/telkomdev/tob/dashboard/utils"
)

// JWTMiddleware this middleware function for verifying accessToken from Authorization Header
func JWTMiddleware(jwtService utils.JwtService, next http.Handler) http.Handler {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		accessToken := req.Header.Get("Authorization")
		if accessToken == "" {
			shared.BuildJSONResponse(res, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    401,
				Message: "no token provided",
				Data:    shared.EmptyJSON{},
			}, 401)
			return
		}
		tokenSlice := strings.Split(accessToken, " ")
		if len(tokenSlice) < 2 {
			shared.BuildJSONResponse(res, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    401,
				Message: "token is not valid",
				Data:    shared.EmptyJSON{},
			}, 401)
			return
		}

		if tokenSlice[0] != "Bearer" {
			shared.BuildJSONResponse(res, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    401,
				Message: "token is not valid",
				Data:    shared.EmptyJSON{},
			}, 401)
			return
		}
		tokenString := tokenSlice[1]
		claim, err := jwtService.Validate(utils.HS256, tokenString)
		if err != nil {
			shared.BuildJSONResponse(res, shared.Response[shared.EmptyJSON]{
				Success: false,
				Code:    401,
				Message: err.Error(),
				Data:    shared.EmptyJSON{},
			}, 401)
			return
		}

		userID := claim.Subject
		req.Header.Add("userId", userID)
		next.ServeHTTP(res, req)
	})
}
