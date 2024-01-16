package mysql

import (
	"database/sql"
	"fmt"
	"gameAppProject/entity"
)

func (d *MySQLDB) IsPhoneNumberUnique(phoneNumber string) (bool, error) {
	user := entity.User{}
	var createdAt []uint8
	// QueryRow executes a query that is expected to return at most one row.
	// If the query selects no rows, the *Row's Scan will return ErrNoRows.
	row := d.db.QueryRow(`select * from users where phone_number =?`, phoneNumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}
		return false, fmt.Errorf("can't scan query result : %w", err)
	}
	return false, nil
}

func (d *MySQLDB) Register(u entity.User) (entity.User, error) {
	res, err := d.db.Exec(`insert into users(name, phone_number , password) values(?,?,?)`, u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("can't execute command %w", err)
	}
	//error is always nil
	id, _ := res.LastInsertId()
	u.ID = uint(id)
	return u, nil
}