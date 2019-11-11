package user

import (
	"database/sql"
	"fmt"
)

type Config struct {
	Username string
	Password string
	Host string
	Database string
	Port int
}

func ConnectDatabase(c Config) (*sql.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
	return sql.Open("mysql", dataSourceName)
}
