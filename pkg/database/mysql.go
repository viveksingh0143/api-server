package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" // Import the MySQL driver
	"github.com/vamika-digital/wms-api-server/config"
)

type MySQLConnection struct {
	Conn *sql.DB
}

func NewMySQLConnection(cfg config.Config) (Connection, error) {
	// Build the connection string from the config
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	)

	// Open a connection to the MySQL database
	conn, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	// Test the connection to ensure it's working
	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &MySQLConnection{Conn: conn}, nil
}
