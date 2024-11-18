package main

import (
	"net/http"

	"github.com/diegopontes87/api/configs"
	"github.com/diegopontes87/api/internal/entity"
	"github.com/diegopontes87/api/internal/infra/database"
	"github.com/diegopontes87/api/internal/infra/webserver/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
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

	http.HandleFunc("/products", productHandler.CreateProduct)
	http.ListenAndServe(":8000", nil)

}
