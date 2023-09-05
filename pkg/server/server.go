package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	authRest "github.com/vamika-digital/wms-api-server/internal/app/auth/interfaces/rest"
	customerRest "github.com/vamika-digital/wms-api-server/internal/app/customer/interfaces"
	machineRest "github.com/vamika-digital/wms-api-server/internal/app/machine/interfaces"
	productRest "github.com/vamika-digital/wms-api-server/internal/app/product/interfaces/rest"
	userRest "github.com/vamika-digital/wms-api-server/internal/app/user/interfaces/rest"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse"
	"github.com/vamika-digital/wms-api-server/internal/middlewares"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type Server struct {
	Address         string
	Port            int
	AuthModule      *authRest.AuthModule
	UserModule      *userRest.UserModule
	ProductModule   *productRest.ProductModule
	MachineModule   *machineRest.MachineModule
	CustomerModule  *customerRest.CustomerModule
	WarehouseModule *warehouse.WarehouseModule
}

func NewServer(address string, port int, db database.Connection) *Server {
	authModule := authRest.NewAuthModule(db)
	userModule := userRest.NewUserModule(db)
	machineModule := machineRest.NewMachineModule(db)
	customerModule := customerRest.NewCustomerModule(db)
	productModule := productRest.NewProductModule(db)
	warehouseModule := warehouse.NewWarehouseModule(db)
	return &Server{
		Address:         address,
		Port:            port,
		AuthModule:      authModule,
		UserModule:      userModule,
		MachineModule:   machineModule,
		CustomerModule:  customerModule,
		ProductModule:   productModule,
		WarehouseModule: warehouseModule,
	}
}

func (s *Server) Run() {
	r := mux.NewRouter()

	r.Use(middlewares.ContentTypeMiddleware)
	r.Use(middlewares.CORSMiddleware)
	s.AuthModule.RegisterRoutes(r.PathPrefix("/auth").Subrouter())
	s.UserModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())
	s.MachineModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())
	s.CustomerModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())
	s.ProductModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())
	s.WarehouseModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())

	log.Printf("Server started on %s:%d", s.Address, s.Port)
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Address, s.Port),
		Handler: r,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server started on %s:%d", s.Address, s.Port)
	}
}
