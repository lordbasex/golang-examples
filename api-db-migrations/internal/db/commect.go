package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ConnectDB crea y retorna una nueva conexión a la base de datos
// Lee la configuración desde variables de entorno
// Uso:
//
//	db := db.ConnectDB()
//	defer db.Close()
func ConnectDB() *sql.DB {
	// Leer variables de entorno
	dsn := getenv("MYSQL_DSN", "api:api@tcp(127.0.0.1:3306)/api?multiStatements=true&parseTime=true")

	log.Printf("[ConnectDB] Intentando conectar a la base de datos...")

	// Intentar conectar con reintentos
	deadline := time.Now().Add(60 * time.Second)
	for {
		db, err := sql.Open("mysql", dsn)
		if err == nil {
			// Configurar el pool de conexiones
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(10)
			db.SetConnMaxLifetime(5 * time.Minute)

			// Validar la conexión con Ping
			err = db.Ping()
			if err == nil {
				log.Printf("[ConnectDB] Conexión establecida exitosamente")
				return db
			}

			log.Printf("[ConnectDB] Fallo al validar conexión: %v. Reintentando...", err)
			db.Close()
		} else {
			log.Printf("[ConnectDB] Fallo al abrir conexión: %v. Reintentando...", err)
		}

		// Verificar timeout
		if time.Now().After(deadline) {
			log.Fatalf("[ConnectDB] No se pudo conectar después de 60 segundos: %v", err)
		}

		time.Sleep(2 * time.Second)
	}
}

// MustConnectDB es como ConnectDB pero hace panic si falla
// Útil para inicialización de la aplicación
func MustConnectDB() *sql.DB {
	db := ConnectDB()
	if db == nil {
		log.Fatal("[MustConnectDB] No se pudo establecer conexión a la base de datos")
	}
	return db
}

// getenv obtiene una variable de entorno con valor por defecto
func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
