// postgres database driver for a bikeshop database
package bikeshop

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func New() {

	uri := "postgres://username:secret@ipAddr:5432/BikeShop"
	os.Setenv("DATABASE_URL", uri)
	ctx := context.Background()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to a database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("This is a msg directly from Bikeshop")
}
