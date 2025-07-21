package user

import (
	"database/sql"
	"log"
	"modernc.org/sqlite"
)

type Model struct {
	Id      int64
	Name    string
	Picture string
}

type RepositoryDriver interface {
	Register(model *Model) error
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

func (repo *Repository) Register(user *Model) error {
	return repo.driver.Register(user)
}
