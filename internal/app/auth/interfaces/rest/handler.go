package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/vamika-digital/wms-api-server/config"
	"github.com/vamika-digital/wms-api-server/internal/app/user/domain"
	"github.com/vamika-digital/wms-api-server/internal/app/user/usecase"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UseCase usecase.UserUseCase
}

func NewAuthHandler(useCase usecase.UserUseCase) *AuthHandler {
	return &AuthHandler{UseCase: useCase}
}

func (handler *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username   string `json:"username"`
		Password   string `json:"password"`
		RememberMe bool   `json:"rememberMe"`
	}

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the credentials
	validUser, err := handler.ValidateCredentials(credentials.Username, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := validUser.GenerateAccessToken(config.AppConfig)
	if err != nil {
		http.Error(w, "Issue at generating access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := validUser.GenerateRefreshToken(config.AppConfig, credentials.RememberMe)
	if err != nil {
		http.Error(w, "Issue at generating refresh token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Name         string `json:"name"`
		StaffID      string `json:"staff_id"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         validUser.Name,
		StaffID:      validUser.StaffID,
	})
}

func (handler *AuthHandler) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	var tokenRequest struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&tokenRequest)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the refresh token
	claims := &jwt.StandardClaims{}
	tkn, err := jwt.ParseWithClaims(tokenRequest.RefreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return config.AppConfig.Auth.SecretKey, nil
	})
	if err != nil || !tkn.Valid {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Get the user
	user, err := handler.UseCase.GetUserByUsername(claims.Subject)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	if user.Status != "active" {
		http.Error(w, "user not active", http.StatusNotFound)
		return
	}

	// Create a new access token
	accessToken, err := user.GenerateAccessToken(config.AppConfig)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// Optionally, create a new refresh token
	refreshToken, err := user.GenerateRefreshToken(config.AppConfig, true)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		Name         string `json:"name"`
		StaffID      string `json:"staff_id"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Name:         user.Name,
		StaffID:      user.StaffID,
	})
}

func (handler *AuthHandler) ValidateCredentials(username string, password string) (*domain.User, error) {
	var user *domain.User
	var err error

	if strings.Contains(username, "@") {
		user, err = handler.UseCase.GetUserByEmail(username)
	} else {
		user, err = handler.UseCase.GetUserByUsername(username)
	}

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, nil
}
