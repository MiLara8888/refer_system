package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"refers_rest/pkg/settings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type RefersDB struct {
	db     *sqlx.DB
	config *settings.Config
	// grpc   *grpc.GrpcSource
}

func New(c *settings.Config) (*RefersDB, error) {

	dsn := c.DB.UrlPostgres()

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		log.Println(err)
		panic(err)
	}

	db.DB.SetMaxOpenConns(20)                        // The default is 0 (unlimited)
	db.DB.SetMaxIdleConns(10)                        // defaultMaxIdleConns = 2
	db.DB.SetConnMaxLifetime(200 * time.Millisecond) // 0, connections are reused forever.
	db.DB.SetConnMaxIdleTime(20 * time.Second)


	schemas := strings.Join(strings.Split(c.DB.Schema, " "), ",")
	_, err = db.ExecContext(context.Background(), fmt.Sprintf(`SET search_path TO "%s"`, strings.ToLower(schemas)))
	if err != nil {
		return nil, err
	}

	DB := &RefersDB{
		db:     db,
		config: c,
	}

	return DB, nil
}

func (s *RefersDB) Close(ctx context.Context) error {
	return s.db.Close()
}
