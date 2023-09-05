package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/machine/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/machine/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/machine/usecase"
	"github.com/vamika-digital/wms-api-server/internal/utility/valueobjects"
)

type MachineHandler struct {
	UseCase usecase.MachineUseCase
}

func NewMachineHandler(useCase usecase.MachineUseCase) *MachineHandler {
	return &MachineHandler{UseCase: useCase}
}

func (handler *MachineHandler) CreateMachine(w http.ResponseWriter, r *http.Request) {
	var machine *domain.Machine = domain.NewMachineWithDefaults()

	if err := json.NewDecoder(r.Body).Decode(machine); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if err := validateMachine(machine); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.CreateMachine(machine); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *MachineHandler) UpdateMachine(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Machine ID", http.StatusBadRequest)
		return
	}

	var machine *domain.Machine = domain.NewMachineWithDefaults()
	if err := json.NewDecoder(r.Body).Decode(machine); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateMachine(machine); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	machine.ID = int64(id)
	if err := handler.UseCase.UpdateMachine(machine); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *MachineHandler) DeleteMachine(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid Machine ID", http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.DeleteMachine(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *MachineHandler) GetMachineByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid Machine ID", http.StatusBadRequest)
		return
	}

	machine, err := handler.UseCase.GetMachineByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(machine); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (handler *MachineHandler) GetAllMachines(w http.ResponseWriter, r *http.Request) {
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
	filterOptions := repository.MachineFilterOptions{
		Name:   r.URL.Query().Get("name"),
		Code:   r.URL.Query().Get("code"),
		Status: r.URL.Query().Get("status"),
	}

	machines, totalMachines, err := handler.UseCase.GetAllMachines(page, pageSize, sort, filterOptions)
	if err != nil {
		// Handle error
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if machines == nil {
		machines = []*domain.Machine{}
	}

	// Create a response with pagination details
	response := valueobjects.PaginatedResponse{
		Data:       machines,
		TotalItems: totalMachines,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: (totalMachines + pageSize - 1) / pageSize, // Calculate total pages
	}

	// Respond with the fetched machines and pagination details
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateMachine(machine *domain.Machine) error {
	if machine.Name == "" {
		return errors.New("name is required")
	}
	if machine.Code == "" {
		return errors.New("code is required")
	}
	if machine.Status == "" {
		return errors.New("status is required")
	}
	return nil
}
