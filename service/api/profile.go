package api

////////////////////example API endpoint
import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"strings"
	"time"

	"github.com/disintegration/imaging"
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
		rt.baseLogger.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !exists {
		if err := rt.db.AddNewUser(username); err != nil {
			rt.baseLogger.Printf("Error adding new user: %s", err.Error())
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
	}
	token, err := GenerateJWT([]byte(secret), username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	respons := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(respons)
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
		_, _ = w.Write([]byte("Username must be between 5 and 16 characters"))
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
		_, _ = w.Write([]byte(err.Error()))
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
	_ = json.NewEncoder(w).Encode(respons)

}

// ###############SetProfilePhoto######################
func (rt *_router) MakePictureFromRequest(r *http.Request) ([]byte, error) {
	// Allow up to 10 MB uploads
	if err := r.ParseMultipartForm(10 << 23); err != nil {
		return nil, errors.New("Error parsing multipart form")
	}

	file, header, err := r.FormFile("photo")
	if err != nil {
		return nil, errors.New("Error retrieving the file from form data")
	}
	// Accept files up to 10 MB
	if header.Size > 10<<23 {
		return nil, errors.New("File size exceeds the 10MB limit")
	}

	//1. sprawdz jakie rozszerzenie
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, errors.New("cannot read file for extension detection")
	}
	_, _ = file.Seek(0, 0)

	mimeType := http.DetectContentType(buffer)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		return nil, errors.New("Only JPG and PNG files are allowed")
	}

	//2. ustaw kwadrat 200x200 px
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, errors.New("cannot decode image")
	}

	img = imaging.Resize(img, 200, 200, imaging.Lanczos)
	//3. return

	buf := new(bytes.Buffer)

	err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return nil, err
	}

	imageBytes := buf.Bytes()

	return imageBytes, nil
}

func (rt *_router) SetProfilePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//sprawdz uprawnienia
	//odbierz zdjecie
	//wysli do MakePictureFromRequest
	//zapisz w DB
	//pusty request wiec zwroc tylko staus

	//rt.baseLogger.Printf("SetProfilePhoto endpoint called")

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		rt.baseLogger.Printf("Error parsing multipart form")
		rt.baseLogger.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uname := r.FormValue("username")
	if uname == "" {
		rt.baseLogger.Printf("Username is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authorised, err := rt.Authorise(w, r, uname)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !authorised {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoData, err := rt.MakePictureFromRequest(r)
	if err != nil {
		rt.baseLogger.Printf(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	err = rt.db.AddProfilePhoto(uname, photoData)
	if err != nil {
		rt.baseLogger.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (rt *_router) GetUserProfilePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	username_photo := ps.ByName("username")
	if username_photo == "" {
		rt.baseLogger.Printf("Username is required")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenStr := authHeader[len("Bearer "):]
	token, err := verifyJWT(tokenStr)
	if err != nil || token == nil || !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	photoData, err := rt.db.GetProfilePhoto(username_photo)
	if err != nil {
		rt.baseLogger.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if photoData == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(photoData)
}
func (rt *_router) SearchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//zwraca listy wszystkich uzytkownikow
	users, err := rt.db.GetAllUsers()
	if err != nil {
		rt.baseLogger.Printf(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(users)
}
