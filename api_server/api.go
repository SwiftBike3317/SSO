package api

import (
	"SSO/storage"
	Types "SSO/types"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

const (
	signingKey = "abracadabra"
)

type CustomClaims struct {
	Userdata Types.Account
	Services []int
	jwt.StandardClaims
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}
func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/login", makeHTTPHandleFunc(s.loginAccount)).Methods("GET")
	router.HandleFunc("/auth", makeHTTPHandleFunc(s.authAccount)).Methods("GET")
	router.HandleFunc("/logout", makeHTTPHandleFunc(s.logoutAccount)).Methods("POST")

	log.Println("JSON API SERVER RUNNING ON PORT: ", s.listenAddr)
	//http.ListenAndServe(s.listenAddr, router)
	srv := &http.Server{Addr: s.listenAddr, Handler: router}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("listenAndServe failed: %v", err)
	}
}

func (s *APIServer) loginAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(Types.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}
	user, exist := s.store.GetUser(req.Email, req.Password)
	if !exist {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "user not exist"})
	}
	services, exist := s.store.GetServices(user.Id)
	if !exist {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "no allowed services"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(3 * time.Hour).Unix(),
		},
		Userdata: user,
		Services: services,
	})

	response, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return WriteJSON(w, http.StatusForbidden, "Failed to sign token")
	}
	WriteJSON(w, http.StatusOK, map[string]interface{}{
		"token": response,
	})

	return nil
}

func (s *APIServer) authAccount(w http.ResponseWriter, r *http.Request) error {

	jwt := r.Header.Get("jwt")
	serviceid, err := strconv.Atoi(r.Header.Get("serviceid"))
	if err != nil {
		return WriteJSON(w, http.StatusForbidden, "Failed to get service id")
	}

	if s.store.InBlacklist(jwt) {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "token in blacklist"})

	}

	claims, err := validateToken(jwt)
	if err != nil {
		return WriteJSON(w, http.StatusForbidden, "Failed to validate token")
	}
	if claims.StandardClaims.ExpiresAt < time.Now().Unix() {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "token expired"})

	}

	if !serviceCheck(serviceid, claims.Services) {
		return WriteJSON(w, http.StatusForbidden, ApiError{Error: "no allowed services"})
	}

	WriteJSON(w, http.StatusOK, claims.Userdata)

	return nil

}

func (s *APIServer) logoutAccount(w http.ResponseWriter, r *http.Request) error {
	jwt := r.Header.Get("jwt")
	s.store.AddToBlacklist(jwt)
	return WriteJSON(w, http.StatusOK, "Token added to blacklist")
}

func validateToken(tokenString string) (CustomClaims, error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})
	if err != nil {
		return CustomClaims{}, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return CustomClaims{}, fmt.Errorf("invalid token")
	}

	return *claims, nil
}
func serviceCheck(serviceID int, userServices []int) bool {
	fmt.Println(serviceID, userServices, "test")
	for _, v := range userServices {
		if v == serviceID {
			return true
		}
	}
	return false
}
