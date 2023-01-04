package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// TokenResponse representa una respuesta JSON que contiene un JWT.
type TokenResponse struct {
	Token string `json:"token"`
}

// secret es una clave privada utilizada para firmar JWTs.
var secret = []byte("my-secret")

// createJWT crea un nuevo JWT para un usuario.
func createJWT(user *User) (string, error) {
	// Crea un nuevo objeto de token, especificando el método de firma y las reclamaciones
	// que deseas que contenga.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	})

	// Firma y obtiene el token codificado completo como una cadena utilizando la clave secreta.
	return token.SignedString(secret)
}

// loginHandler maneja las solicitudes de inicio de sesión.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Parsea el cuerpo de la solicitud para obtener las credenciales de inicio de sesión del usuario.
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Payload de solicitud no válido", http.StatusBadRequest)
		return
	}

	// Verifica si las credenciales de inicio de sesión proporcionadas son válidas.
	dbUser, err := verifyCredentials(user.Username, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Crea un nuevo JWT para el usuario.
	token, err := createJWT(dbUser)
	if err != nil {
		http.Error(w, "Error al crear el token", http.StatusInternalServerError)
		return
	}

	// Devuelve el token al cliente.
	json.NewEncoder(w).Encode(TokenResponse{Token: token})
}
