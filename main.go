package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func main() {
	conn, err := pgx.Connect(context.Background(),
		"postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	fmt.Println("Connected to DB!")

	insertJob(conn)
}

func insertJob(conn *pgx.Conn) {
	query := `
	INSERT INTO jobs (payload)
	VALUES ($1)
	RETURNING public_id;
	`

	payload := `{
		"original_filename": "test.jpg",
		"stored_path": "uploads/test.jpg",
		"mime_type": "image/jpeg",
		"file_size": 123456
	}`

	var publicID string

	err := conn.QueryRow(context.Background(), query, payload).Scan(&publicID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New Job Public ID:", publicID)
}

func getNextJob(conn *pgx.Conn) {
	query := `
	SELECT id, payload
	FROM jobs
	WHERE status = 'pending'
	ORDER BY created_at
	LIMIT 1;
	`

	var id string
	var payload string

	err := conn.QueryRow(context.Background(), query).Scan(&id, &payload)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Next job:", id, payload)
}

func updateProgress(conn *pgx.Conn, id string) {
	query := `
	UPDATE jobs
	SET progress = 50,
	    result = result || '{"blurred_path":"uploads/blur.jpg"}'::jsonb
	WHERE id = $1;
	`

	_, err := conn.Exec(context.Background(), query, id)
	if err != nil {
		log.Fatal(err)
	}
}

func checkStatus(conn *pgx.Conn, publicID string) {
	query := `
	SELECT status, progress, result, error_msg
	FROM jobs
	WHERE public_id = $1;
	`

	var status string
	var progress int
	var result string
	var errorMsg *string

	err := conn.QueryRow(context.Background(), query, publicID).
		Scan(&status, &progress, &result, &errorMsg)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Status:", status)
	fmt.Println("Progress:", progress)
	fmt.Println("Result:", result)
}