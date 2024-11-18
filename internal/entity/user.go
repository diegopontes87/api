package entity

import (
	"github.com/diegopontes87/api/pkg/service"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       service.ID `json:"id"`
	Name     string     `json:"name"`
	Password string     `json:"-"`
	Email    string     `json:"email"`
}

func NewUser(name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:       service.NewID(),
		Name:     name,
		Password: string(hash),
		Email:    email,
	}, nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
