package data

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GetDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/landscape.db")
	if err != nil {
		fmt.Println(err)
	}
	return db
}

// Struct representing the data

// Function to open the SQLite database
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Function to create the table if not exists
func CreateTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS registries (
		id TEXT PRIMARY KEY,
		name TEXT,
		logo TEXT,
		category TEXT,
		subcategory TEXT
	)`)
	if err != nil {
		return err
	}
	return nil
}

// Function to insert data into the database
func InsertData(db *sql.DB, data RegistryItem) error {
    // insert if not exists
    _, err := db.Exec("INSERT OR IGNORE INTO registries (id, name, logo, category, subcategory) VALUES (?, ?, ?, ?, ?)",
		data.ID, data.Name, data.Logo, data.Category, data.Subcategory)
	if err != nil {
		return err
	}
	return nil
}

// Function to update data in the database
func UpdateData(db *sql.DB, data RegistryItem) error {
	_, err := db.Exec("UPDATE registries SET name=?, logo=?, category=?, subcategory=? WHERE id=?",
		data.Name, data.Logo, data.Category, data.Subcategory, data.ID)
	if err != nil {
		return err
	}
	return nil
}

// Function to delete data from the database
func DeleteData(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM registries WHERE id=?", id)
	if err != nil {
		return err
	}
	return nil
}

func GetRegistry(db *sql.DB) ([]RegistryItem, error) {
	var registries []RegistryItem
	rows, err := db.Query("SELECT id, name, logo, category, subcategory FROM registries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id, name, logo, category, subcategory string
		if err := rows.Scan(&id, &name, &logo, &category, &subcategory); err != nil {
			return nil, err
		}
		registries = append(registries, RegistryItem{ID: id, Name: name, Logo: logo, Category: category, Subcategory: subcategory})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return registries, nil
}
