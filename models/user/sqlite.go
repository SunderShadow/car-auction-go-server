package user

import (
	"database/sql"
	"errors"
)

type sqliteDriver struct {
	db *sql.DB
}

func (driver sqliteDriver) FindById(id int64) *Model {
	row := driver.db.QueryRow(`SELECT * FROM users WHERE id = ?`, id)

	model := new(Model)

	err := row.Scan(&model.Id, &model.Name, &model.Picture)

	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return model
}

func (driver sqliteDriver) Register(model *Model) error {
	res, resErr := driver.db.Exec("INSERT INTO users(name, picture) VALUES (?, ?)", model.Name, model.Picture)

	if resErr != nil {
		return resErr
	}

	id, idErr := res.LastInsertId()

	if idErr != nil {
		return idErr
	}

	model.Id = id

	return nil
}

func (driver sqliteDriver) CreateTable() error {
	_, err := driver.db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name VARCHAR,
		    picture varchar
		)
	`)

	return err
}
