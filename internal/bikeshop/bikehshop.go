// postgres database driver for a bikeshop database
package bikeshop

import (
	"context"
	"fmt"
	"os"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Invoice struct {
	Fname    string  `json:"fname"`
	Lname    string  `json:"lname"`
	Product  string  `json:"product"`
	Price    float32 `json:"price"`
	Quantity int     `json:"quantity"`
	Category string  `json:"category"`
	Shipping string  `json:"shipping"`
}

type invoices []*Invoice

func (invs invoices) string() string {
	str := ""
	for _, inv := range invs {
		str += fmt.Sprintf("%v\n", *inv)
	}
	return str
}

func New() {

	uri := "postgres://username@localhost:5432/bikeshop"
	os.Setenv("DATABASE_URL", uri)
	ctx := context.Background()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to a database: %v\n", err)
		os.Exit(1)
	}

	defer db.Close()

	var invs invoices
	if err := pgxscan.Select(ctx, db, &invs, `SELECT fname, lname,
        product, price, quantity FROM invoices`); err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
	}

	fmt.Printf("All the invoices\n")
	fmt.Println(invs)
}
