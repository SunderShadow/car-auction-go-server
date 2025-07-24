package lot

import (
	"database/sql"
	"log"
	"modernc.org/sqlite"
)

type Model struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Picture     string `json:"picture"`
	Description string `json:"description"`
	CurrentBid  int64  `json:"current_bid"`
}

type RepositoryDriver interface {
	Create(model *Model) error
	FindAll() []*Model
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

func (repo *Repository) CreateTable() error {
	return repo.driver.CreateTable()
}

func (repo *Repository) FindAll() []*Model {
	return repo.driver.FindAll()
}

func (repo *Repository) Create(model *Model) error {
	return repo.driver.Create(model)
}
