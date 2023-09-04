package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/usecase"
	"github.com/vamika-digital/wms-api-server/internal/utility/valueobjects"
)

type InventoryHandler struct {
	UseCase usecase.InventoryUseCase
}

func NewInventoryHandler(useCase usecase.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{UseCase: useCase}
}

// subRouter.HandleFunc("/raw-material", u.Handler.CreateInventoryForRawMaterial).Methods(http.MethodPost)
// subRouter.HandleFunc("/finished-goods", u.Handler.CreateInventoryForFinishedGoods).Methods(http.MethodPost)

func (handler *InventoryHandler) CreateInventoryForRawMaterial(w http.ResponseWriter, r *http.Request) {
	var inventoryForm *domain.InventoryFormRawMaterial = &domain.InventoryFormRawMaterial{}

	if err := json.NewDecoder(r.Body).Decode(inventoryForm); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateInventoryFormRawMaterial(inventoryForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.CreateInventoryForRawMaterial(inventoryForm); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *InventoryHandler) CreateInventory(w http.ResponseWriter, r *http.Request) {
	var inventory *domain.Inventory = domain.NewInventoryWithDefaults()

	if err := json.NewDecoder(r.Body).Decode(inventory); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateInventory(inventory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.CreateInventory(inventory); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *InventoryHandler) UpdateInventory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
		return
	}

	var inventory *domain.Inventory = domain.NewInventoryWithDefaults()
	if err := json.NewDecoder(r.Body).Decode(inventory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateInventory(inventory); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	inventory.ID = int64(id)
	if err := handler.UseCase.UpdateInventory(inventory); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *InventoryHandler) DeleteInventory(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid Inventory ID", http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.DeleteInventory(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *InventoryHandler) GetInventoryByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid Inventory ID", http.StatusBadRequest)
		return
	}

	inventory, err := handler.UseCase.GetInventoryByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(inventory); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *InventoryHandler) GetAllInventories(w http.ResponseWriter, r *http.Request) {
	pageQuery := r.URL.Query().Get("page")
	pageSizeQuery := r.URL.Query().Get("pageSize")
	sort := r.URL.Query().Get("sort")

	// Default values for pagination
	page := 1
	pageSize := 10

	if pageQuery != "" {
		page, _ = strconv.Atoi(pageQuery)
	}
	if pageSizeQuery != "" {
		pageSize, _ = strconv.Atoi(pageSizeQuery)
	}

	// Extract filter parameters
	filterOptions := repository.InventoryFilterOptions{
		// Name: r.URL.Query().Get("name"),
		// Status:    r.URL.Query().Get("status"),
		// ProductID: r.URL.Query().Get("name"),
		// RackID:    r.URL.Query().Get("name"),
		// StoreID:   r.URL.Query().Get("name"),
	}
	// filterOptions.SetOwnerID(r.URL.Query().Get("owner_id"))

	inventories, totalInventories, err := handler.UseCase.GetAllInventories(page, pageSize, sort, filterOptions)
	if err != nil {
		// Handle error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if inventories == nil {
		inventories = []*domain.Inventory{}
	}

	// Create a response with pagination details
	response := valueobjects.PaginatedResponse{
		Data:       inventories,
		TotalItems: totalInventories,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (totalInventories + pageSize - 1) / pageSize, // Calculate total pages
	}

	// Respond with the fetched inventories and pagination details
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateInventory(inventory *domain.Inventory) error {
	if inventory.Status == "" {
		return errors.New("status is required")
	}
	return nil
}

func validateInventoryFormRawMaterial(inventory *domain.InventoryFormRawMaterial) error {
	if inventory.Product_id == "" {
		return errors.New("product is required")
	}
	if inventory.Pallet == "" {
		return errors.New("pallet code is required")
	}
	if inventory.Quantity <= 0 {
		return errors.New("quantity is required")
	}
	return nil
}
