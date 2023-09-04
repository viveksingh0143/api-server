package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/usecase"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type InventoryModule struct {
	Handler *InventoryHandler
}

func NewInventoryModule(db database.Connection) *InventoryModule {
	inventoryRepo := repository.NewInventoryRepository(db)
	inventoryUsecase := usecase.NewInventoryUseCase(inventoryRepo)
	inventoryHandler := NewInventoryHandler(inventoryUsecase)

	return &InventoryModule{Handler: inventoryHandler}
}

func (u *InventoryModule) RegisterRoutes(r *mux.Router) {
	subRouter := r.PathPrefix("/inventories").Subrouter()

	subRouter.HandleFunc("/raw-material", u.Handler.CreateInventoryForRawMaterial).Methods(http.MethodPost)
	// subRouter.HandleFunc("/finished-goods", u.Handler.CreateInventoryForFinishedGoods).Methods(http.MethodPost)

	subRouter.HandleFunc("", u.Handler.CreateInventory).Methods(http.MethodPost)
	subRouter.HandleFunc("", u.Handler.GetAllInventories).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.GetInventoryByID).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.UpdateInventory).Methods(http.MethodPut)
	subRouter.HandleFunc("/{id}", u.Handler.DeleteInventory).Methods(http.MethodDelete)
}
