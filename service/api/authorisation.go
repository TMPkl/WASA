package api

import (
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func (rt *_router) Authorise(w http.ResponseWriter, r *http.Request, username string) (bool, error) {
	authHeader := r.Header.Get("Authorization")
	if exist, _ := rt.db.UserExists(username); !exist {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("User does not exist"))
		return false, errors.New("User does not exist")
	}
	if authHeader == "" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authorization header is missing"))
		return false, errors.New("Authorization header is missing")
	}
	tokenStr := authHeader[len("Bearer "):]
	token, err := verifyJWT(tokenStr)
	if err != nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token"))
		return false, errors.New("Invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token claims"))
		return false, errors.New("Invalid token claims")
	}
	username_token, ok := claims["sub"].(string)
	if !ok || username_token == "" || username_token != username {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid token subject"))
		return false, errors.New("Invalid token subject")
	}
	return true, nil
}
