package main

import (
	"net/http"

	"github.com/diegopontes87/api/configs"
	"github.com/diegopontes87/api/internal/entity"
	"github.com/diegopontes87/api/internal/infra/database"
	"github.com/diegopontes87/api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName     string = "test.db"
	configPath string = "../../configs"
)

func main() {

	cfg, err := configs.LoadConfig(configPath)
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUserDB(db)
	userHandler := handlers.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", cfg.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", cfg.JWTExpiresIn))

	r.Route("/products", func(r chi.Router) {
		r.Use(middleware.WithValue("jwt", cfg.TokenAuth))
		r.Use(jwtauth.Verifier(cfg.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Post("/generate_token", userHandler.GetJWT)
	})

	http.ListenAndServe(":8000", r)
}
