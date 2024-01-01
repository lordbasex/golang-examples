package database

import (
	"api-fiber/config"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

const tableUsers = `CREATE TABLE IF NOT EXISTS users (
	id int(11) NOT NULL AUTO_INCREMENT,
	fullname varchar(50) NOT NULL,
	email varchar(100) NOT NULL,
	pass varchar(50) NOT NULL DEFAULT '',
	creation_date datetime NOT NULL,
	PRIMARY KEY (id),
	UNIQUE KEY email (email),
	KEY ct_users_email (email)
)`

const tableUsersKeys = `CREATE TABLE IF NOT EXISTS users_keys (
    id int(11) NOT NULL AUTO_INCREMENT,
	user_id int(11) NOT NULL,
    pvkey text,
    pbkey text,
	creation_date datetime NOT NULL,
	PRIMARY KEY (id),
	KEY users_keys_user_id (user_id),
	CONSTRAINT ct_users_id FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
)`

// DBInit initializes database tables when the daemon starts.
func DBInit() {
	// Connect to the database
	db := ConnectD()

	// Create the 'users' table
	_, errCreateTableUser := db.Exec(tableUsers)
	if errCreateTableUser != nil {
		// Log error if debug mode is enabled
		if config.DebugConfiguration {
			fmt.Print(errCreateTableUser)
		}
		return
	}

	// Create the 'users_keys' table
	_, errCreateTableUserKey := db.Exec(tableUsersKeys)
	if errCreateTableUserKey != nil {
		// Log error if debug mode is enabled
		if config.DebugConfiguration {
			fmt.Print(errCreateTableUserKey)
		}
		return
	}
}
