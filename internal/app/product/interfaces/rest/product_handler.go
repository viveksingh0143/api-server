package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/product/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/product/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/product/usecase"
	"github.com/vamika-digital/wms-api-server/internal/utility/valueobjects"
)

type ProductHandler struct {
	UseCase usecase.ProductUseCase
}

func NewProductHandler(useCase usecase.ProductUseCase) *ProductHandler {
	return &ProductHandler{UseCase: useCase}
}

func (handler *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product domain.Product = domain.NewProductWithDefaults()

	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateProduct(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.CreateProduct(&product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	var product domain.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateProduct(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	product.ID = int64(id)
	if err := handler.UseCase.UpdateProduct(&product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid Product ID", http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.DeleteProduct(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid Product ID", http.StatusBadRequest)
		return
	}

	product, err := handler.UseCase.GetProductByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(product); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
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
	filterOptions := repository.ProductFilterOptions{
		Type:    r.URL.Query().Get("type"),
		Code:    r.URL.Query().Get("code"),
		RawCode: r.URL.Query().Get("raw_code"),
		Name:    r.URL.Query().Get("name"),
		Status:  r.URL.Query().Get("status"),
	}

	products, totalProducts, err := handler.UseCase.GetAllProducts(page, pageSize, sort, filterOptions)
	if err != nil {
		// Handle error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if products == nil {
		products = []*domain.Product{}
	}

	// Create a response with pagination details
	response := valueobjects.PaginatedResponse{
		Data:       products,
		TotalItems: totalProducts,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (totalProducts + pageSize - 1) / pageSize, // Calculate total pages
	}

	// Respond with the fetched products and pagination details
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateProduct(product *domain.Product) error {
	typeErr := product.ValidateType()
	if typeErr != nil {
		return errors.New("type should be valid")
	}
	if product.Code == "" {
		return errors.New("code is required")
	}
	if product.Name == "" {
		return errors.New("name is required")
	}
	if product.Description == "" {
		return errors.New("description is description")
	}
	if product.Unit == "" {
		return errors.New("unit is required")
	}
	if product.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
