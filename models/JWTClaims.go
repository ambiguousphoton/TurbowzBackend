package models

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	UserID     int64    `json:"user_id"`
	UserHandle string `json:"user_handle"`
	jwt.RegisteredClaims
}