package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
)

func Initialize() {
	//db, err := sql.Open("postgres", "user=postgres password=password sslmode=disable")
	dburl := "postgres://postgres:password@localhost:5432/imageboard?sslmode=disable"
	dbpool, err := pgxpool.New(context.Background(), dburl)
	defer dbpool.Close()

	if err != nil {
		log.Fatal(err)
	}

	_, err = dbpool.Exec(context.Background(), "CREATE DATABASE imageboard")
	if err != nil {
		log.Println(err)
	}

	log.Println("Database successfully initialized")

	db := Establish()
	CreateTables(db)

	defer db.Close()
}

func CreateTables(db *pgxpool.Pool) {

	file, _ := os.ReadFile("database/util/create.sql")
	_, err := db.Exec(context.Background(), string(file))
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Table posts successfully initialized")
	}

}

func Establish() *pgxpool.Pool {

	dburl := "postgres://postgres:password@localhost:5432/imageboard?sslmode=disable"
	dbpool, err := pgxpool.New(context.Background(), dburl)
	if err != nil {
		log.Fatal(err)
	}
	err = dbpool.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return dbpool
}
