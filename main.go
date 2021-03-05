package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Employee struct {
	Id     int    `json:"id"`
	Gender string `json:"gender"`
}

func employees(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, gender FROM employees")
		if err != nil {
			fmt.Fprintf(w, "server error")
			return
		}
		defer rows.Close()
		var employees []Employee
		for rows.Next() {
			var e Employee
			if err := rows.Scan(&e.Id, &e.Gender); err != nil {
				fmt.Fprintf(w, "server error")
				return
			}
			employees = append(employees, e)
		}
		resp, err := json.Marshal(employees)
		if err != nil {
			fmt.Fprintf(w, "server error")
			return
		}
		w.Write(resp)
	}
}

func initDB(dbPath string) (*sql.DB, error) {
	if info, err := os.Stat(dbPath); os.IsNotExist(err) {
		file, err := os.Create(dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create DB file: %s", err)
		}
		file.Close()
	} else if info.IsDir() {
		return nil, fmt.Errorf("db_path %q points to directory, instead of file", dbPath)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB file: %s", err)
	}
	return db, nil
}

func createDbSchema(db *sql.DB) error {
	employees := `
        CREATE TABLE IF NOT EXISTS employees (
            id INT PRIMARY KEY,        
            gender TEXT
        );`
	if _, err := db.Exec(employees); err != nil {
		return fmt.Errorf("failed to create employees table: %s", err)
	}
	return nil
}

func seedDb(db *sql.DB) error {
	insert := `
        DELETE FROM employees;
        INSERT INTO employees
            (id, gender)
        VALUES 
            (1, "male"),
            (2, "male"),
            (3, "male"),
            (4, "female"),
            (5, "female"),
            (6, "female");`
	if _, err := db.Exec(insert); err != nil {
		return fmt.Errorf("failed to seed employees table: %s", err)
	}
	return nil
}

func main() {
	var port int
	var dbPath string
	var seed bool
	port, _ = strconv.Atoi(os.Getenv("PORT"))
	flag.IntVar(&port, "port", port, "Port number for serving HTTP requests")
	if port < 0 || port > 65353 {
		panic(fmt.Sprintf("invalid port number, expecting between 0 and 65353, but got %d", port))
	}
	flag.StringVar(&dbPath, "db_path", "./employees.db", "Path to SQLite DB file")
	flag.BoolVar(&seed, "seed_db", false, "Seed DB with initial employee data")
	flag.Parse()
	db, err := initDB(dbPath)
	if err != nil {
		panic(fmt.Sprintf("failed to open DB file: %s", err))
	}
	if err := createDbSchema(db); err != nil {
		panic(fmt.Sprintf("failed to create DB schema: %s", err))
	}
	if seed {
		if err := seedDb(db); err != nil {
			panic(fmt.Sprintf("failed to seed DB %s", err))
		}
	}
	http.HandleFunc("/employees", employees(db))
	fmt.Printf("Serving employee data from %q on localhost:%d\n", dbPath, port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
