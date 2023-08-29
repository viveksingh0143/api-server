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

type StoreHandler struct {
	UseCase usecase.StoreUseCase
}

func NewStoreHandler(useCase usecase.StoreUseCase) *StoreHandler {
	return &StoreHandler{UseCase: useCase}
}

func (handler *StoreHandler) CreateStore(w http.ResponseWriter, r *http.Request) {
	var store *domain.Store = domain.NewStoreWithDefaults()

	if err := json.NewDecoder(r.Body).Decode(store); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateStore(store); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.CreateStore(store); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *StoreHandler) UpdateStore(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Store ID", http.StatusBadRequest)
		return
	}

	var store *domain.Store = domain.NewStoreWithDefaults()
	if err := json.NewDecoder(r.Body).Decode(store); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateStore(store); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	store.ID = int64(id)
	if err := handler.UseCase.UpdateStore(store); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *StoreHandler) DeleteStore(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid Store ID", http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.DeleteStore(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *StoreHandler) GetStoreByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid Store ID", http.StatusBadRequest)
		return
	}

	store, err := handler.UseCase.GetStoreByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(store); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *StoreHandler) GetAllStores(w http.ResponseWriter, r *http.Request) {
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
	filterOptions := repository.StoreFilterOptions{
		Name:     r.URL.Query().Get("name"),
		Location: r.URL.Query().Get("location"),
		Status:   r.URL.Query().Get("status"),
	}
	filterOptions.SetOwnerID(r.URL.Query().Get("owner_id"))

	stores, totalStores, err := handler.UseCase.GetAllStores(page, pageSize, sort, filterOptions)
	if err != nil {
		// Handle error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if stores == nil {
		stores = []*domain.Store{}
	}

	// Create a response with pagination details
	response := valueobjects.PaginatedResponse{
		Data:       stores,
		TotalItems: totalStores,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (totalStores + pageSize - 1) / pageSize, // Calculate total pages
	}

	// Respond with the fetched stores and pagination details
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateStore(store *domain.Store) error {
	if store.Name == "" {
		return errors.New("name is required")
	}
	if store.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
