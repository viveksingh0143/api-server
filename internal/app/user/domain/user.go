package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vamika-digital/wms-api-server/config"
)

type User struct {
	ID                int64     `json:"id"`
	Username          string    `json:"username"`
	Password          string    `json:"password"`
	PasswordHash      string    `json:"-"`
	Name              string    `json:"name"`
	StaffID           string    `json:"staff_id"`
	Email             string    `json:"email"`
	EmailConfirmation bool      `json:"email_confirmation"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	LastUpdatedBy     string    `json:"last_updated_by"`
}

func NewUserWithDefaults() User {
	return User{
		EmailConfirmation: false,
		Status:            "active",
	}
}

func (u *User) GenerateAccessToken(cfg config.Config) (string, error) {
	expirationTime := time.Now().Add(time.Second * time.Duration(cfg.Auth.ExpiryDuration)).Unix()
	claims := &jwt.StandardClaims{
		Subject:   u.Username,
		ExpiresAt: expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.SecretKey))
}

func (u *User) GenerateRefreshToken(cfg config.Config, expireLong bool) (string, error) {
	var expirationTime int64
	if expireLong {
		expirationTime = time.Now().Add(time.Hour * 24 * time.Duration(cfg.Auth.ExpiryLongDuration)).Unix()
	} else {
		expirationTime = time.Now().Add(time.Second * time.Duration(cfg.Auth.ExpiryDuration)).Unix()
	}

	claims := &jwt.StandardClaims{
		Subject:   u.Username,
		ExpiresAt: expirationTime,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.SecretKey))
}
