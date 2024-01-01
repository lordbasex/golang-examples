package models

import "github.com/golang-jwt/jwt"

// User represents the structure of a user.
type User struct {
	FullName string `json:"fullname"` // Full name of the user
	Email    string `json:"email"`    // Email address of the user
	Password string `json:"password"` // Password of the user
}

// CustomClaims represents the structure of custom claims in the token.
type CustomClaims struct {
	UserInfo           User `json:"user"` // User information embedded in the token claims
	jwt.StandardClaims      // Standard JWT claims (expiration time, issuance time, etc.)
}
