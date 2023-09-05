package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/customer/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/customer/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/customer/usecase"
	"github.com/vamika-digital/wms-api-server/internal/utility/valueobjects"
)

type CustomerHandler struct {
	UseCase usecase.CustomerUseCase
}

func NewCustomerHandler(useCase usecase.CustomerUseCase) *CustomerHandler {
	return &CustomerHandler{UseCase: useCase}
}

func (handler *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer *domain.Customer = domain.NewCustomerWithDefaults()

	if err := json.NewDecoder(r.Body).Decode(customer); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateCustomer(customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.CreateCustomer(customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Customer ID", http.StatusBadRequest)
		return
	}

	var customer *domain.Customer = domain.NewCustomerWithDefaults()
	if err := json.NewDecoder(r.Body).Decode(customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateCustomer(customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	customer.ID = int64(id)
	if err := handler.UseCase.UpdateCustomer(customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid Customer ID", http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.DeleteCustomer(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *CustomerHandler) GetCustomerByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid Customer ID", http.StatusBadRequest)
		return
	}

	customer, err := handler.UseCase.GetCustomerByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *CustomerHandler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
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
	filterOptions := repository.CustomerFilterOptions{
		Name:   r.URL.Query().Get("name"),
		Code:   r.URL.Query().Get("code"),
		Status: r.URL.Query().Get("status"),
	}

	customers, totalCustomers, err := handler.UseCase.GetAllCustomers(page, pageSize, sort, filterOptions)
	if err != nil {
		// Handle error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if customers == nil {
		customers = []*domain.Customer{}
	}

	// Create a response with pagination details
	response := valueobjects.PaginatedResponse{
		Data:       customers,
		TotalItems: totalCustomers,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalCustomers / pageSize, // Calculate total pages
	}

	// Respond with the fetched customers and pagination details
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateCustomer(customer *domain.Customer) error {
	if customer.Name == "" {
		return errors.New("name is required")
	}
	if customer.Code == "" {
		return errors.New("code is required")
	}
	if customer.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
