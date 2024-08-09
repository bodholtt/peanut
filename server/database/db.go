package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"peanutserver/pcfg"
	"strconv"
)

func Initialize() {
	//db, err := sql.Open("postgres", "user=postgres password=password sslmode=disable")
	//dburl := "postgres://postgres:password@localhost:5432/imageboard?sslmode=disable"
	dburl := "postgres://" +
		pcfg.Cfg.Database.Username + ":" +
		pcfg.Cfg.Database.Password + "@" +
		pcfg.Cfg.Database.Host + ":" +
		strconv.Itoa(pcfg.Cfg.Database.Port) + "/" +
		pcfg.Cfg.Database.DatabaseName +
		pcfg.Cfg.Database.Params

	dbpool, err := pgxpool.New(context.Background(), dburl)
	defer dbpool.Close()

	if err != nil {
		log.Fatal(err)
	}

	_, err = dbpool.Exec(context.Background(), "CREATE DATABASE "+pcfg.Cfg.Database.DatabaseName)
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

// Establish - Establish a database connection
// pool - Connection pool *** MUST BE CLOSED ***
func Establish() (pool *pgxpool.Pool) {

	dburl := "postgres://" +
		pcfg.Cfg.Database.Username + ":" +
		pcfg.Cfg.Database.Password + "@" +
		pcfg.Cfg.Database.Host + ":" +
		strconv.Itoa(pcfg.Cfg.Database.Port) + "/" +
		pcfg.Cfg.Database.DatabaseName +
		pcfg.Cfg.Database.Params

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
