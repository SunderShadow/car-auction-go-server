package google

import (
	"car-auction/models/user"
	"database/sql"
	"log"
	"modernc.org/sqlite"
)

type Model struct {
	UserId               int64
	AccessToken          string
	AccessTokenExpiresIn int64
	GoogleUserId         string
}

type RepositoryDriver interface {
	FindByGoogleUserId(id string) *Model
	Register(user *user.Model, model *Model) error
	CreateTable() error
}

type Repository struct {
	driver RepositoryDriver
}

func NewRepository(db *sql.DB) *Repository {
	switch db.Driver().(type) {
	case *sqlite.Driver:
		return &Repository{
			sqliteDriver{db},
		}
	default:
		log.Fatalln(`No repository driver for current db driver`)
		return nil
	}
}

func (repo *Repository) FindByGoogleUserId(id string) *Model {
	return repo.driver.FindByGoogleUserId(id)
}

func (repo *Repository) CreateTable() error {
	return repo.driver.CreateTable()
}

func (repo *Repository) Register(user *user.Model, model *Model) error {
	return repo.driver.Register(user, model)
}
