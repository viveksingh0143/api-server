package warehouse

import (
	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/warehouse/interfaces/rest"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type WarehouseModule struct {
	ContainerModule *rest.ContainerModule
	StoreModule     *rest.StoreModule
}

func NewWarehouseModule(db database.Connection) *WarehouseModule {
	containerModule := rest.NewContainerModule(db)
	storeModule := rest.NewStoreModule(db)
	return &WarehouseModule{ContainerModule: containerModule, StoreModule: storeModule}
}

func (w *WarehouseModule) RegisterRoutes(r *mux.Router) {
	w.ContainerModule.RegisterRoutes(r.PathPrefix("/warehouse").Subrouter())
	w.StoreModule.RegisterRoutes(r.PathPrefix("/warehouse").Subrouter())
}
