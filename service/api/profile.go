package api

////////////////////example API endpoint
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

// /login request struct -> onl username in body
type LoginRequest struct {
	Username string `json:"username"`
}

var secret = "my_not_so_secret_key"

// checkUser sprawdza czy użytkownik istnieje w bazie danych
func (rt *_router) checkUser(r *http.Request) (string, bool, error) {
	var req LoginRequest

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&req)

	if err != nil {
		rt.baseLogger.Printf(err.Error())
		return "", false, errors.New("Invalid request body")
	}

	if req.Username == "" {
		err := errors.New("Username is required")
		return "", false, err
	}

	// czy uzytkownik juz istnieje
	exists, err := rt.db.UserExists(req.Username)
	if err != nil {
		return "", false, err
	}

	return req.Username, exists, nil
}

func GenerateJWT(secret []byte, username string) (string, error) {
	claims := jwt.MapClaims{
		"sub": username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func verifyJWT(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})
}

// /login endpoint
func (rt *_router) login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username, exists, err := rt.checkUser(r)
	if err != nil {
		rt.baseLogger.Printf(string(err.Error()))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !exists {
		rt.db.AddNewUser(username)
	}
	token, err := GenerateJWT([]byte(secret), username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	respons := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respons)
}

// ##############UpdateMyUsername##################
type UpdateMyUsernameRequest struct {
	Username    string `json:"username"`
	NewUsername string `json:"new-username"`
}

func (rt *_router) UpdateMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req UpdateMyUsernameRequest
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&req)

	if len(req.NewUsername) < 5 || len(req.NewUsername) > 16 || req.NewUsername == "" {
		rt.baseLogger.Printf("Username must be between 5 and 16 characters")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username must be between 5 and 16 characters"))
		return
	}

	if err != nil {
		rt.baseLogger.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		rt.baseLogger.Printf("Invalid request body")
		return
	}

	authorised, err := rt.Authorise(w, r, req.Username)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authorised {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = rt.db.UpdateUsername(req.Username, req.NewUsername)
	if err != nil {
		rt.baseLogger.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	token, err := GenerateJWT([]byte(secret), req.NewUsername)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	respons := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respons)

}
