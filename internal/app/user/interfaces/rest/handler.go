package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/user/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/user/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/user/usecase"
	"github.com/vamika-digital/wms-api-server/internal/utility/valueobjects"
)

type UserHandler struct {
	UseCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{UseCase: useCase}
}

func (handler *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User = domain.NewUserWithDefaults()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.CreateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = int64(id)
	if err := handler.UseCase.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid User ID", http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.DeleteUser(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	user, err := handler.UseCase.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
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
	filterOptions := repository.UserFilterOptions{
		Name:     r.URL.Query().Get("name"),
		Username: r.URL.Query().Get("username"),
		Email:    r.URL.Query().Get("email"),
		Status:   r.URL.Query().Get("status"),
	}

	users, totalUsers, err := handler.UseCase.GetAllUsers(page, pageSize, sort, filterOptions)
	if err != nil {
		// Handle error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create a response with pagination details
	response := valueobjects.PaginatedResponse{
		Data:       users,
		TotalItems: totalUsers,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (totalUsers + pageSize - 1) / pageSize, // Calculate total pages
	}

	// Respond with the fetched users and pagination details
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateUser(user *domain.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}

	if user.Name == "" {
		return errors.New("name is required")
	}

	if user.Email == "" {
		return errors.New("email is required")
	}

	if user.Password == "" {
		return errors.New("password is required")
	}

	// You can add further validation checks here...

	return nil
}
