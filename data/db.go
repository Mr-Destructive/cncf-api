package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func loadEnv() map[string]string {
	// open .env file and load envs
	file_path := ".env"
	// read file
	bytes, err := os.ReadFile(file_path)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	envs := map[string]string{}
	for _, line := range strings.Split(string(bytes), "\n") {
		pair := strings.SplitN(line, "=", 2)
		if len(pair) == 2 {
			envs[pair[0]] = pair[1]
		}
	}
	return envs
}

func GetDb() *sql.DB {
	envs := loadEnv()
	dbName := envs["TURSO_DB_NAME"]
	orgName := envs["TURSO_ORG_NAME"]
	authToken := envs["TURSO_API_TOKEN"]
	url := fmt.Sprintf("libsql://%s-%s.turso.io?authToken=%s", dbName, orgName, authToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db %s: %s", url, err)
		os.Exit(1)
	}
	return db
}

/*
func GetDb() *sql.DB {
	db, err := sql.Open("sqlite3", "./data/landscape.db")
	if err != nil {
		fmt.Println(err)
	}
	return db
}
*/

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
	);`)
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

func GetRegistry(db *sql.DB, filter string, args ...interface{}) ([]RegistryItem, error) {
	var registries []RegistryItem
	var rows *sql.Rows
	var err error

	query := "SELECT id, name, logo, category, subcategory FROM registries"
	if filter != "" {
		query += " WHERE " + filter
	}

	rows, err = db.Query(query, args...)
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
