package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/diegopontes87/api/internal/dto"
	"github.com/diegopontes87/api/internal/entity"
	"github.com/diegopontes87/api/internal/infra/database"
	"github.com/diegopontes87/api/pkg/service"
	"github.com/go-chi/chi"
)

type ProductHandler struct {
	ProductDB database.ProductDBInterface
}

func NewProductHandler(db database.ProductDBInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// CreateProduct godoc
// @Summary      Create a product
// @Description  Create a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request  body   dto.CreateProductInput  true  "product request"
// @Success      201
// @Failure      400     {object}  entity.Error
// @Failure      500     {object}  entity.Error
// @Router       /products [post]
// @Security ApiKeyAuth
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product *dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetProduct	 godoc
// @Summary      Get product
// @Description  Get product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id  path        string     true "product ID" Format(uuid)
// @Success      200   			 {object}   entity.Product
// @Failure      400     	     {object}  entity.Error
// @Failure      404
// @Router       /products/{id}	 [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := entity.Error{Message: "ID cant be nil"}
		json.NewEncoder(w).Encode(err)
		return
	}
	product, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(product)
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path   string		             true  "product ID" Format(uuid)
// @Param        request  body   dto.CreateProductInput  true  "product request"
// @Success      200
// @Failure      404
// @Failure      400     {object}  entity.Error
// @Failure      500     {object}  entity.Error
// @Router       /products [put]
// @Security ApiKeyAuth
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := entity.Error{Message: "ID cant be nil"}
		json.NewEncoder(w).Encode(err)
		return
	}
	var product entity.Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	product.ID, err = service.ParseID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	_, err = h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.ProductDB.Update(&product)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)

}

// DeleteProduct godoc
// @Summary      Delete a product
// @Description  Delete a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path     string		    true  "product ID" Format(uuid)
// @Success      200
// @Failure      404
// @Failure      400     {object}  entity.Error
// @Failure      500     {object}  entity.Error
// @Router       /products/{id}    [delete]
// @Security ApiKeyAuth
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := entity.Error{Message: "ID cant be nil"}
		json.NewEncoder(w).Encode(err)
	}
	_, err := h.ProductDB.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = h.ProductDB.Delete(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetProducts 	 godoc
// @Summary      List products
// @Description  get all products
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request page 	 query	 	string false "page number"
// @Param        request limit   query	 	string false "limit"
// @Success      200   			 {array}    entity.Product
// @Failure      500   			 {object}   entity.Error
// @Router       /products 		 [get]
// @Security ApiKeyAuth
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	limit := r.URL.Query().Get("limit")
	sort := r.URL.Query().Get("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 0
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		limitInt = 0
	}

	products, err := h.ProductDB.FindAll(pageInt, limitInt, sort)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err := entity.Error{Message: err.Error()}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&products)
}
