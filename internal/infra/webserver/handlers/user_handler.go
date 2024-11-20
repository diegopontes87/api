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

// CreateUser    godoc
// @Summary      Create a user
// @Description  Create a new user in the application with the provided data
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body   dto.CreateUserInput  true  "User request"
// @Success      201
// @Failure      500   {object}  entity.Error
// @Failure      400   {object}  entity.Error
// @Router       /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.UserDB.Create(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetJWT    	 godoc
// @Summary      Get a user JWT
// @Description  Get a user JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request  body   dto.GetJWTInput  true  "user credentials"
// @Success      200   {object}  dto.GetJWTOutput
// @Failure      404   {object}  entity.Error
// @Failure      500   {object}  entity.Error
// @Router       /users/generate_token [post]
func (h *UserHandler) GetJWT(w http.ResponseWriter, r *http.Request) {
	var userJWT dto.GetJWTInput
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("JwtExpiresIn").(int)

	err := json.NewDecoder(r.Body).Decode(&userJWT)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
	}
	user, err := h.UserDB.FindByEmail(userJWT.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
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

	accessToken := dto.GetJWTOutput{AccessToken: tokenString}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(accessToken)
}
