package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var UserValues []interface{}

// jwtMiddleware verifica el JWT en las solicitudes que no sean las de inicio de sesión.
func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			response := &Msg{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Unauthorized",
			}
			responseWriter(http.StatusBadRequest, response, w)
			return
		}

		// Valida el formato del encabezado (Bearer seguido por un espacio y el token)
		if !strings.HasPrefix(tokenString, "Bearer ") {
			if DebugConfiguration {
				log.Printf("Auth Header Error: incorrect format")
			}

			response := &Msg{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Unauthorized",
			}
			responseWriter(http.StatusBadRequest, response, w)
			return
		}

		parts := strings.Split(tokenString, "Bearer")
		tokenNew := strings.TrimSpace(parts[1])

		if len(tokenNew) < 1 {
			if DebugConfiguration {
				log.Printf("Token Header Error Bearer: %s", parts)
			}

			response := &Msg{
				Code:    http.StatusUnauthorized,
				Status:  "error",
				Message: "Unauthorized",
			}
			responseWriter(http.StatusBadRequest, response, w)
			return
		}

		if DebugConfiguration {
			log.Printf("authHeader: %s", tokenNew)
		}

		token, err := jwt.Parse(tokenNew, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				response := &Msg{
					Code:    http.StatusUnauthorized,
					Status:  "error",
					Message: "unexpected_signature_method",
				}
				responseWriter(http.StatusBadRequest, response, w)
			}
			return []byte(Secret), nil
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			// Iterar a través del mapa y almacenar los valores en el slice
			for _, value := range claims {
				UserValues = append(UserValues, value)
			}

			// Asignar los valores a las posiciones correctas en UserValues
			for key, value := range claims {
				switch key {
				case "id":
					UserValues[0] = value
				case "user":
					UserValues[1] = value
				case "exp":
					UserValues[2] = value
				}
			}

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

			response := &Msg{
				Code:    http.StatusMethodNotAllowed,
				Status:  "error",
				Message: "unexpected_signature_method",
			}
			responseWriter(http.StatusMethodNotAllowed, response, w)
		}
		// Llama a la función siguiente en la cadena de middlewares/handlers
		next(w, r)
	}
}

// authMiddleware es una función de middleware que aplica la función de middleware
// corsMiddleware y methodMiddleware a la función que maneja la ruta.
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	if DebugConfiguration {
		log.Print("func: authMiddleware")
	}
	// La función de middleware devuelve una nueva función http.HandlerFunc
	return corsMiddleware(methodMiddleware(next))
}
