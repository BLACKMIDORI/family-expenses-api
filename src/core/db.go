package core

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
)

func getConnection() (*pgx.Conn, error) {
	ctx := context.Background()
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbUrl := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf(
		"postgresql://%v:%v@%v/%v?sslmode=verify-full",
		dbUser,
		dbPassword,
		dbUrl,
		dbName,
	)
	conn, err := pgx.Connect(ctx, connectionString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	conn, err := getConnection()
	if err != nil {
		return nil, err
	}
	transaction, err := conn.BeginTx(ctx, pgx.TxOptions{})
	return transaction, err
}
