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

type ContainerHandler struct {
	UseCase usecase.ContainerUseCase
}

func NewContainerHandler(useCase usecase.ContainerUseCase) *ContainerHandler {
	return &ContainerHandler{UseCase: useCase}
}

func (handler *ContainerHandler) CreateContainer(w http.ResponseWriter, r *http.Request) {
	var container domain.Container = domain.NewContainerWithDefaults()

	if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateContainer(&container); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.CreateContainer(&container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *ContainerHandler) UpdateContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Container ID", http.StatusBadRequest)
		return
	}

	var container domain.Container
	if err := json.NewDecoder(r.Body).Decode(&container); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateContainer(&container); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	container.ID = int64(id)
	if err := handler.UseCase.UpdateContainer(&container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ContainerHandler) DeleteContainer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid Container ID", http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.DeleteContainer(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *ContainerHandler) GetContainerByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid Container ID", http.StatusBadRequest)
		return
	}

	container, err := handler.UseCase.GetContainerByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(container); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *ContainerHandler) GetAllContainers(w http.ResponseWriter, r *http.Request) {
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
	filterOptions := repository.ContainerFilterOptions{
		Type:   r.URL.Query().Get("type"),
		Code:   r.URL.Query().Get("code"),
		Name:   r.URL.Query().Get("name"),
		Status: r.URL.Query().Get("status"),
	}

	containers, totalContainers, err := handler.UseCase.GetAllContainers(page, pageSize, sort, filterOptions)
	if err != nil {
		// Handle error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if containers == nil {
		containers = []*domain.Container{}
	}

	// Create a response with pagination details
	response := valueobjects.PaginatedResponse{
		Data:       containers,
		TotalItems: totalContainers,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (totalContainers + pageSize - 1) / pageSize, // Calculate total pages
	}

	// Respond with the fetched containers and pagination details
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateContainer(container *domain.Container) error {
	typeErr := container.ValidateType()
	if typeErr != nil {
		return errors.New("type should be valid")
	}
	if container.Code == "" {
		return errors.New("code is required")
	}
	if container.Name == "" {
		return errors.New("name is required")
	}
	if container.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
