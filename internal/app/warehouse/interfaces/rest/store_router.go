package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/usecase"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type StoreModule struct {
	Handler *StoreHandler
}

func NewStoreModule(db database.Connection) *StoreModule {
	storeRepo := repository.NewStoreRepository(db)
	storeUsecase := usecase.NewStoreUseCase(storeRepo)
	storeHandler := NewStoreHandler(storeUsecase)

	return &StoreModule{Handler: storeHandler}
}

func (u *StoreModule) RegisterRoutes(r *mux.Router) {
	subRouter := r.PathPrefix("/stores").Subrouter()
	subRouter.HandleFunc("", u.Handler.CreateStore).Methods(http.MethodPost)
	subRouter.HandleFunc("", u.Handler.GetAllStores).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.GetStoreByID).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.UpdateStore).Methods(http.MethodPut)
	subRouter.HandleFunc("/{id}", u.Handler.DeleteStore).Methods(http.MethodDelete)
}
