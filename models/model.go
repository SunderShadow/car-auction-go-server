package models

type DBModel interface {
	DBCreateTable() error
}
