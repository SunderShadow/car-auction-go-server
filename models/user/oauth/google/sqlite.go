package google

import (
	"car-auction/models/user"
	"database/sql"
)

type sqliteDriver struct {
	db *sql.DB
}

func (driver sqliteDriver) CreateTable() error {
	_, err := driver.db.Exec(`
		CREATE TABLE IF NOT EXISTS google_oauth(
		    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		    access_token VARCHAR,
		    access_token_expires_in INTEGER,
		    
		    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)

	return err
}

func (driver sqliteDriver) Register(user *user.Model, model *Model) error {
	_, resErr := driver.db.Exec("INSERT INTO google_oauth(user_id, access_token, access_token_expires_in) VALUES (?, ?, ?)", user.Id, model.AccessToken, model.AccessTokenExpiresIn)

	if resErr != nil {
		return resErr
	}

	model.UserId = user.Id

	return nil
}
