package google

import (
	"car-auction/models/user"
	"database/sql"
	"errors"
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
		    google_user_id VARCHAR,
		    current_bid INTEGER,
		    
		    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)

	return err
}

func (driver sqliteDriver) FindByGoogleUserId(id string) *Model {
	row := driver.db.QueryRow("SELECT * FROM google_oauth WHERE google_user_id = ?", id)

	model := new(Model)

	err := row.Scan(&model.UserId, &model.AccessToken, &model.AccessTokenExpiresIn, &model.GoogleUserId)

	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}

	return model
}

func (driver sqliteDriver) Register(user *user.Model, model *Model) error {
	_, resErr := driver.db.Exec(`INSERT INTO 
    	google_oauth(user_id, access_token, access_token_expires_in, google_user_id) 
		VALUES (?, ?, ?, ?)
	`, user.Id, model.AccessToken, model.AccessTokenExpiresIn, model.GoogleUserId)

	if resErr != nil {
		return resErr
	}

	model.UserId = user.Id

	return nil
}
