package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/user/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/user/usecase"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type UserModule struct {
	Handler *UserHandler
}

func NewUserModule(db database.Connection) *UserModule {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepo)
	userHandler := NewUserHandler(userUsecase)

	return &UserModule{Handler: userHandler}
}

func (u *UserModule) RegisterRoutes(r *mux.Router) {
	subRouter := r.PathPrefix("/users").Subrouter()

	subRouter.HandleFunc("", u.Handler.CreateUser).Methods(http.MethodPost)
	subRouter.HandleFunc("", u.Handler.GetAllUsers).Methods(http.MethodGet)
	subRouter.HandleFunc("/{id}", u.Handler.GetUserByID).Methods(http.MethodGet)
	subRouter.HandleFunc("/{id}", u.Handler.UpdateUser).Methods(http.MethodPut)
	subRouter.HandleFunc("/{id}", u.Handler.DeleteUser).Methods(http.MethodDelete)
}
