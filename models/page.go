package models

import (
	"database/sql"
	"time"
)

type Page struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS pages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT UNIQUE NOT NULL,
		content TEXT NOT NULL,
		description TEXT,
		author TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (p *Page) Create(db *sql.DB) error {
	query := `INSERT INTO pages (title, slug, content, description, author) VALUES (?, ?, ?, ?, ?)`
	result, err := db.Exec(query, p.Title, p.Slug, p.Content, p.Description, p.Author)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

func (p *Page) Update(db *sql.DB) error {
	query := `UPDATE pages SET title = ?, content = ?, description = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.Exec(query, p.Title, p.Content, p.Description, p.ID)
	return err
}

func GetPage(db *sql.DB, slug string) (*Page, error) {
	query := `SELECT id, title, slug, content, description, author, created_at, updated_at FROM pages WHERE slug = ?`
	row := db.QueryRow(query, slug)

	var page Page
	err := row.Scan(&page.ID, &page.Title, &page.Slug, &page.Content, &page.Description, &page.Author, &page.CreatedAt, &page.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &page, nil
}

func GetAllPages(db *sql.DB) ([]Page, error) {
	query := `SELECT id, title, slug, content, description, author, created_at, updated_at FROM pages ORDER BY created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pages []Page
	for rows.Next() {
		var page Page
		err := rows.Scan(&page.ID, &page.Title, &page.Slug, &page.Content, &page.Description, &page.Author, &page.CreatedAt, &page.UpdatedAt)
		if err != nil {
			return nil, err
		}
		pages = append(pages, page)
	}

	return pages, nil
}

func DeletePage(db *sql.DB, id int) error {
	query := `DELETE FROM pages WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
