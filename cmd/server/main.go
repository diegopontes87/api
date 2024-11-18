package main

import (
	"net/http"

	"github.com/diegopontes87/api/configs"
	"github.com/diegopontes87/api/internal/entity"
	"github.com/diegopontes87/api/internal/infra/database"
	"github.com/diegopontes87/api/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	dbName     string = "test.db"
	configPath string = "../../configs"
)

func main() {

	_, err := configs.LoadConfig(configPath)
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

	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products", productHandler.GetProducts)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)

	r.Post("/users", userHandler.CreateUser)

	http.ListenAndServe(":8000", r)
}
