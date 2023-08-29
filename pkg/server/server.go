package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	authRest "github.com/vamika-digital/wms-api-server/internal/app/auth/interfaces/rest"
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
	WarehouseModule *warehouse.WarehouseModule
}

func NewServer(address string, port int, db database.Connection) *Server {
	authModule := authRest.NewAuthModule(db)
	userModule := userRest.NewUserModule(db)
	productModule := productRest.NewProductModule(db)
	warehouseModule := warehouse.NewWarehouseModule(db)
	return &Server{Address: address, Port: port, AuthModule: authModule, UserModule: userModule, ProductModule: productModule, WarehouseModule: warehouseModule}
}

func (s *Server) Run() {
	r := mux.NewRouter()

	r.Use(middlewares.ContentTypeMiddleware)
	r.Use(middlewares.CORSMiddleware)
	s.AuthModule.RegisterRoutes(r.PathPrefix("/auth").Subrouter())
	s.UserModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())
	s.ProductModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())
	s.WarehouseModule.RegisterRoutes(r.PathPrefix("/secure").Subrouter())

	// c := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:3000"}, // Your allowed origin
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	// 	AllowedHeaders:   []string{"Authorization", "Content-Type"}, // You can add any required headers here
	// 	AllowCredentials: true,
	// })

	// Wrap the router with the CORS handler
	// handler := c.Handler(r)

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
