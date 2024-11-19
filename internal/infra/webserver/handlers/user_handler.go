package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/diegopontes87/api/internal/dto"
	"github.com/diegopontes87/api/internal/entity"
	"github.com/diegopontes87/api/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB database.UserDBInterface
}

func NewUserHandler(db database.UserDBInterface) *UserHandler {
	return &UserHandler{
		UserDB: db,
	}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var userJWT dto.GetJWTInput
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)

	err := json.NewDecoder(r.Body).Decode(&userJWT)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	user, err := h.UserDB.FindByEmail(userJWT.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !user.ValidatePassword(userJWT.Password) {
		fmt.Printf("Password is not valid: %v", userJWT.Password)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, tokenString, _ := jwt.Encode(map[string]interface{}{
		"sub": user.ID.String(),
		"exp": time.Now().Add(time.Second * time.Duration(jwtExpiresIn)).Unix(),
	})

	accessToken := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: tokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}
