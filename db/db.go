package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/Tutuacs/pkg/config"
)

type Store struct {
}

func NewConnection() (conn *sql.DB, err error) {
	conf := config.GetDB()

	stringConnection := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Pass, conf.Name)

	conn, err = sql.Open("postgres", stringConnection)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()

	return
}

func (s *Store) CloseConnection(conn *sql.DB) error {
	return conn.Close()
}
