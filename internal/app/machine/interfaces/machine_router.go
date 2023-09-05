package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/machine/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/machine/usecase"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type MachineModule struct {
	Handler *MachineHandler
}

func NewMachineModule(db database.Connection) *MachineModule {
	machinesRepo := repository.NewMachineRepository(db)
	machinesUsecase := usecase.NewMachineUseCase(machinesRepo)
	machinesHandler := NewMachineHandler(machinesUsecase)

	return &MachineModule{Handler: machinesHandler}
}

func (u *MachineModule) RegisterRoutes(r *mux.Router) {
	subRouter := r.PathPrefix("/machines").Subrouter()
	subRouter.HandleFunc("", u.Handler.CreateMachine).Methods(http.MethodPost)
	subRouter.HandleFunc("", u.Handler.GetAllMachines).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.GetMachineByID).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.UpdateMachine).Methods(http.MethodPut)
	subRouter.HandleFunc("/{id}", u.Handler.DeleteMachine).Methods(http.MethodDelete)
}
