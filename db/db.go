package db

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	db, err := sql.Open("sqlite3", makeConnString())
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	if err := setupDB(db); err != nil {
		panic(err)
	}

	Get = bindGetDB(db)
}

var Get func() *sql.DB

func bindGetDB(db *sql.DB) func() *sql.DB {
	return func() *sql.DB {
		return db
	}
}

func setupDB(db *sql.DB) error {
	sets := [...]string{
		`CREATE TABLE IF NOT EXISTS settings (
          id INTEGER PRIMARY KEY,  
          payload JSON NOT NULL
    )`,
	}

	tx, err := db.Begin()
	for _, q := range sets {
		_, err = db.Exec(q)
	}

	if err != nil {
		tx.Rollback()
	} else {
		tx.Commit()
	}
	return err
}

func makeConnString() string {
	return fmt.Sprintf("file:telegospel.db")
}

func SelectSettingsFor(id int64) (data map[string]any, err error) {
	db := Get()
	res := db.QueryRow(`
    SELECT payload 
    FROM settings 
    WHERE id = $1;
    `, id,
	)

	if err = res.Err(); err != nil {
		return
	}

	var jsonLoad []byte
	res.Scan(&jsonLoad)

	err = json.Unmarshal(jsonLoad, &data)
	if err != nil {
		return
	}

	return
}

func InsertSettingsFor(id int64, data map[string]any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	db := Get()

	_, err = db.Exec(`
    INSERT INTO settings(id, payload)
    VALUES($1, $2)
    ON CONFLICT DO UPDATE
    SET payload = $2;
    `, id, payload,
	)

	return err
}
