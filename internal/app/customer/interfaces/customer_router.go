package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/customer/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/customer/usecase"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type CustomerModule struct {
	Handler *CustomerHandler
}

func NewCustomerModule(db database.Connection) *CustomerModule {
	customersRepo := repository.NewCustomerRepository(db)
	customersUsecase := usecase.NewCustomerUseCase(customersRepo)
	customersHandler := NewCustomerHandler(customersUsecase)

	return &CustomerModule{Handler: customersHandler}
}

func (u *CustomerModule) RegisterRoutes(r *mux.Router) {
	subRouter := r.PathPrefix("/customers").Subrouter()
	subRouter.HandleFunc("", u.Handler.CreateCustomer).Methods(http.MethodPost)
	subRouter.HandleFunc("", u.Handler.GetAllCustomers).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.GetCustomerByID).Methods(http.MethodGet, http.MethodOptions)
	subRouter.HandleFunc("/{id}", u.Handler.UpdateCustomer).Methods(http.MethodPut)
	subRouter.HandleFunc("/{id}", u.Handler.DeleteCustomer).Methods(http.MethodDelete)
}
