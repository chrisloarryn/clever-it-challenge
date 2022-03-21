package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"

	"CleverIT-challenge/internal/core/domain/beers"
)

type Repository struct {
	*sql.DB
}

func (repository *Repository) FindAllBeers(ctx context.Context) ([]beers.Beer, error) {
	queryStmt := `select "id", "name", "brewery", "country", "currency", "price" from "beer"`
	rows, err := repository.DB.QueryContext(ctx, queryStmt)
	if err != nil {
		return nil, err
	}
	var result []beers.Beer
	for rows.Next() {
		rowResult := beers.Beer{}
		if err := rows.Scan(&rowResult.ID, &rowResult.Name, &rowResult.Brewery, &rowResult.Country, &rowResult.Currency, &rowResult.Price); err != nil {
			return nil, err
		}
		result = append(result, rowResult)
	}
	return result, nil
}

func (repository *Repository) FindBeerByID(ctx context.Context, beerID int) (beers.Beer, error) {
	queryStmt := `select "id", "name", "brewery", "country", "currency", "price" from "beer" where ID = $1`
	row := repository.DB.QueryRowContext(ctx, queryStmt, beerID)
	if row == nil {
		return beers.Beer{}, fmt.Errorf("no result for ID: %d", beerID)
	}
	result := beers.Beer{}
	if err := row.Scan(&result.ID, &result.Name, &result.Brewery, &result.Country, &result.Currency, &result.Price); err != nil {
		return beers.Beer{}, err
	}
	return result, nil
}

func (repository *Repository) SaveBeer(ctx context.Context, beer beers.Beer) error {
	insertDynStmt := `insert into "beer"("id", "name", "brewery", "country", "currency", "price") values($1, $2, $3, $4, $5, $6)`
	_, err := repository.DB.ExecContext(ctx, insertDynStmt, beer.ID, beer.Name, beer.Brewery, beer.Country, beer.Currency, beer.Price)
	if err != nil {
		return err
	}
	return nil
}

var (
	dbHostKey     = "DB_HOST"
	dbPortKey     = "DB_PORT"
	dbDatabaseKey = "DB_DATABASE"
	dbUserKey     = "DB_USER"
	dbPasswordKey = "DB_PASSWORD"
)

func NewRepository() beers.Repository {
	host := os.Getenv(dbHostKey)
	if len(host) == 0 {
		host = "localhost"
	}
	var err error
	portStr := os.Getenv(dbPortKey)
	port := 5432
	if len(portStr) > 0 {
		if portInt, err := strconv.Atoi(portStr); err != nil {
			port = portInt
		}
	}
	var user = os.Getenv(dbUserKey)
	if len(user) == 0 {
		user = "postgres"
	}
	var password = os.Getenv(dbPasswordKey)
	if len(password) == 0 {
		password = "root"
	}
	var dbname = os.Getenv(dbDatabaseKey)
	if len(dbname) == 0 {
		dbname = "postgres"
	}
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", sqlConn)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// check db
	err = db.Ping()
	if err != nil {
		log.Fatalln(err.Error())
	}

	log.Println("Connected!")

	return &Repository{
		db,
	}
}
