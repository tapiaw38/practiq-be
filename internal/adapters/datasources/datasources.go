package datasources

import "database/sql"

type Datasources struct {
	DB *sql.DB
}

func CreateDatasources(db *sql.DB) *Datasources {
	return &Datasources{
		DB: db,
	}
}
