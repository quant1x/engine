package storages

import (
	"database/sql"
	"errors"
	"fmt"
	"gitee.com/quant1x/engine/config"
	_ "github.com/marcboeker/go-duckdb"
	"log"
)

type Database struct {
	db *sql.DB
}

func NewDatabaseStorage(db *sql.DB) *Database {
	return &Database{db: db}
}

func (s *Database) Create() {
	_, err := s.db.Exec(`CREATE TABLE people (id INTEGER, name VARCHAR)`)
	if err != nil {
		log.Println(err)
	}
	_, err = s.db.Exec(`INSERT INTO people VALUES (42, 'John')`)
	if err != nil {
		log.Println(err)
	}
}

func (s *Database) Query() {
	var (
		id   int
		name string
	)
	row := s.db.QueryRow(`SELECT id, name FROM people`)
	err := row.Scan(&id, &name)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println("no rows")
	} else if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("id: %d, name: %s\n", id, name)
}

func Connect() *sql.DB {
	quant1XConfig, found := config.LoadConfig()
	if !found {
		log.Fatal("Failed to find config file")
	}

	db, err := sql.Open("duckdb", quant1XConfig.BaseDir+"/duck.db")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
