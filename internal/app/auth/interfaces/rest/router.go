package rest

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/app/user/repository"
	"github.com/vamika-digital/wms-api-server/internal/app/user/usecase"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type AuthModule struct {
	Handler *AuthHandler
}

func NewAuthModule(db database.Connection) *AuthModule {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepo)
	AuthHandler := NewAuthHandler(userUsecase)

	return &AuthModule{Handler: AuthHandler}
}

func (u *AuthModule) RegisterRoutes(r *mux.Router) {
	subRouter := r

	subRouter.HandleFunc("/login", u.Handler.LoginHandler).Methods(http.MethodPost, http.MethodOptions)
	subRouter.HandleFunc("/refresh-token", u.Handler.RefreshTokenHandler).Methods(http.MethodPost, http.MethodOptions)
}
