package db

import (
	"database/sql"
	"fmt"
	"reflect"

	_ "github.com/lib/pq"

	"github.com/Tutuacs/pkg/config"
)

type Store struct {
}

// type StoreInterface interface {
// 	NewConnection() (*sql.DB, error)
// 	ScanRow(*sql.Row, interface{}) error
// 	ScanRows(*sql.Rows, interface{}) error
// }

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

func ScanRow(row *sql.Row, targetType interface{}) (err error) {

	val := reflect.ValueOf(targetType).Elem()
	numFields := val.NumField()
	scanArgs := make([]interface{}, numFields)
	for i := 0; i < numFields; i++ {
		field := val.Field(i)
		scanArgs[i] = field.Addr().Interface()
	}

	err = row.Scan(scanArgs...)

	return
}

func ScanRows(rows *sql.Rows, targetType interface{}) (err error) {

	val := reflect.ValueOf(targetType).Elem()

	numFields := val.NumField()

	scanArgs := make([]interface{}, numFields)
	for i := 0; i < numFields; i++ {
		field := val.Field(i)
		scanArgs[i] = field.Addr().Interface()
	}

	err = rows.Scan(scanArgs...)

	return
}
