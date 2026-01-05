package api

////////////////////example API endpoint
import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

var secret = "my_not_so_secret_key"

// checkUser sprawdza czy użytkownik istnieje w bazie danych
func (rt *_router) checkUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) (string, bool, error) {
	// Pobierz username z query parameters
	username := r.URL.Query().Get("username")
	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.New("Username query parameter is required")
		return "", false, err
	}

	// Sprawdź czy użytkownik istnieje
	exists, err := rt.db.UserExists(username)
	if err != nil {
		return "", false, err
	}

	return username, exists, nil
}

func GenerateJWT(secret []byte, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func VerifyJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secret, nil
	})
}

func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}
