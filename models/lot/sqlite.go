package lot

import (
	"database/sql"
)

type sqliteDriver struct {
	db *sql.DB
}

func (driver sqliteDriver) Create(model *Model) error {
	res, err := driver.db.Exec(
		`INSERT INTO auction_lots(name, description, picture, current_bid) VALUES(?, ?, ?, ?)`,
		model.Name, model.Description, model.Picture, model.CurrentBid,
	)

	if err != nil {
		return err
	}

	model.Id, _ = res.LastInsertId()

	return nil
}

func (driver sqliteDriver) FindAll() []*Model {
	rows, _ := driver.db.Query("SELECT * FROM auction_lots")

	models := make([]*Model, 0)

	for rows.Next() {
		var model Model
		rows.Scan(&model.Id, &model.Name, &model.Description, &model.Picture, &model.CurrentBid)

		models = append(models, &model)
	}

	return models
}

func (driver sqliteDriver) CreateTable() error {
	_, err := driver.db.Exec(`
		CREATE TABLE IF NOT EXISTS auction_lots(
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name VARCHAR,
		    description TEXT,
		    picture TEXT,
		    current_bid INTEGER
		)
	`)

	return err
}
