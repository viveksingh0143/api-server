package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/product/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/product/usecase"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type ProductModule struct {
	Handler *ProductHandler
}

func NewProductModule(db database.Connection) *ProductModule {
	productRepo := repository.NewProductRepository(db)
	productUsecase := usecase.NewProductUseCase(productRepo)
	productHandler := NewProductHandler(productUsecase)

	return &ProductModule{Handler: productHandler}
}

func (u *ProductModule) RegisterRoutes(r *mux.Router) {
	subRouter := r.PathPrefix("/products").Subrouter()
	subRouter.HandleFunc("", u.Handler.CreateProduct).Methods(http.MethodPost)
	subRouter.HandleFunc("", u.Handler.GetAllProducts).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.GetProductByID).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.UpdateProduct).Methods(http.MethodPut)
	subRouter.HandleFunc("/{id}", u.Handler.DeleteProduct).Methods(http.MethodDelete)
}
