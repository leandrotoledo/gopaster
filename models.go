package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	db *sql.DB
)

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "data/gopaster.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS pastes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		password TEXT NOT NULL
    )`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatal(err)
	}
}

type Paste struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	Password  string
}

func (p *Paste) Save() error {
	var query string
	var args []interface{}

	if p.Password != "" {
		hashedPassword, err := hashPassword(p.Password)
		if err != nil {
			return err
		}
		query = "INSERT INTO pastes (title, content, created_at, password) VALUES (?, ?, CURRENT_TIMESTAMP, ?)"
		args = append(args, p.Title, p.Content, hashedPassword)
	} else {
		query = "INSERT INTO pastes (title, content, created_at, password) VALUES (?, ?, CURRENT_TIMESTAMP, '')"
		args = append(args, p.Title, p.Content)
	}

	res, err := db.Exec(query, args...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	return p.Load(int(id))
}

func (p *Paste) Load(id int) error {
	query := "SELECT id, title, content, created_at, password FROM pastes WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&p.ID, &p.Title, &p.Content, &p.CreatedAt, &p.Password)

	return err
}

func (p *Paste) Delete() error {
	query := "DELETE FROM pastes WHERE id = ?"
	_, err := db.Exec(query, p.ID)
	return err
}

func loadPasteByID(id int) (*Paste, error) {
	paste := &Paste{}
	if err := paste.Load(id); err != nil {
		return nil, err
	}
	return paste, nil
}

func loadAllPastes() ([]*Paste, error) {
	query := "SELECT id, title, content, created_at FROM pastes ORDER BY created_at DESC"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pastes []*Paste
	for rows.Next() {
		paste := &Paste{}
		if err := rows.Scan(&paste.ID, &paste.Title, &paste.Content, &paste.CreatedAt); err != nil {
			return nil, err
		}
		pastes = append(pastes, paste)
	}
	return pastes, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
