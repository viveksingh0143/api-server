package database

import "database/sql"

type Connection interface {
	GetDB() *sql.DB
}

func (m *MySQLConnection) GetDB() *sql.DB {
	return m.Conn
}
