package main

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

// jwtMiddleware verifica el JWT en las solicitudes que no sean las de inicio de sesión.
func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Falta el token JWT", http.StatusUnauthorized)
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Printf("%v\n", claims)
		} else {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// corsMiddleware es una función de middleware que agrega la cabecera
// Access-Control-Allow-Origin: * a todas las respuestas HTTP. Esto permite
// el acceso a los recursos del servidor desde cualquier dominio.
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// La función de middleware devuelve una nueva función http.HandlerFunc
	return func(w http.ResponseWriter, r *http.Request) {

		// Agrega el user-agent a la respuesta HTTP
		w.Header().Set("User-Agent", r.Header.Get("User-Agent"))

		// Agrega la cabecera Access-Control-Allow-Origin: * a la respuesta HTTP
		w.Header().Set("Access-Control-Allow-Origin", "*")
		// Llama a la función siguiente en la cadena de middlewares/handlers
		next(w, r)
	}
}

// methodMiddleware es una función de middleware que bloquea métodos HTTP que no
// sean GET, POST, PUT o DELETE.
func methodMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// La función de middleware devuelve una nueva función http.HandlerFunc
	return func(w http.ResponseWriter, r *http.Request) {
		// Si el método de la solicitud no es permitido, devuelve una respuesta HTTP
		// con el código de estado 405 Method Not Allowed y termina la ejecución
		// de la función de middleware.
		if r.Method != "GET" && r.Method != "POST" && r.Method != "PUT" && r.Method != "DELETE" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		// Llama a la función siguiente en la cadena de middlewares/handlers
		next(w, r)
	}
}

// authMiddleware es una función de middleware que aplica la función de middleware
// corsMiddleware y methodMiddleware a la función que maneja la ruta.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	// La función de middleware devuelve una nueva función http.HandlerFunc
	return corsMiddleware(methodMiddleware(next))
}
