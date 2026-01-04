package db

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	mmysql "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MustRunMigrations aplica migraciones UP. No requiere agregar columnas ni tocar tus tablas.
// Solo usa una tabla interna para controlar versiones (schema_migrations).
func MustRunMigrations(db *sql.DB, migrationsPath string) {
	driver, err := mmysql.WithInstance(db, &mmysql.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Printf("migrations: no changes")
			return
		}
		log.Fatal(err)
	}
	log.Printf("migrations: applied successfully")
}
