package guards

import (
	"database/sql"
	"fmt"

	"github.com/Tutuacs/pkg/db"
	"github.com/Tutuacs/pkg/types"
)

type Store struct {
	db *sql.DB
}

func (s *Store) GetUserByEmail(email string) (usr *types.User, err error) {
	err = nil
	usr = &types.User{}

	query := "SELECT * FROM  users WHERE email = $1"
	row := s.db.QueryRow(query, email)

	db.ScanRow(row, usr)

	if usr.ID == 0 {
		err = fmt.Errorf("user not found")
		return
	}

	return
}

func (s *Store) GetUserByID(ID int) (*types.User, error) {

	sql := "SELECT * FROM users WHERE id = $1"

	rows, err := s.db.Query(sql, ID)
	if err != nil {
		return nil, err
	}

	usr := new(types.User)

	for rows.Next() {
		err = db.ScanRows(rows, usr)
		if err != nil {
			return nil, err
		}
	}

	return usr, err
}
