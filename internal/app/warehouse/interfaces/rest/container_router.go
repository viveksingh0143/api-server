package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/usecase"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type ContainerModule struct {
	Handler *ContainerHandler
}

func NewContainerModule(db database.Connection) *ContainerModule {
	containerRepo := repository.NewContainerRepository(db)
	containerUsecase := usecase.NewContainerUseCase(containerRepo)
	containerHandler := NewContainerHandler(containerUsecase)

	return &ContainerModule{Handler: containerHandler}
}

func (u *ContainerModule) RegisterRoutes(r *mux.Router) {
	subRouter := r.PathPrefix("/containers").Subrouter()
	subRouter.HandleFunc("", u.Handler.CreateContainer).Methods(http.MethodPost)
	subRouter.HandleFunc("", u.Handler.GetAllContainers).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.GetContainerByID).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.UpdateContainer).Methods(http.MethodPut)
	subRouter.HandleFunc("/{id}", u.Handler.DeleteContainer).Methods(http.MethodDelete)
}
