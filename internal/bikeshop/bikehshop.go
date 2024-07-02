package bikeshop

import (
	"context"
	"fmt"
	"os"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Invoice struct {
	Fname    string  `form:"fname"`
	Lname    string  `form:"lname"`
	Product  string  `form:"product"`
	Price    float32 `form:"price"`
	Quantity int     `form:"quantity"`
	Category string  `form:"category"`
	Shipping string  `form:"shipping"`
}

type invoices []*Invoice

// implements stringer interface to print each invoice
func (invs invoices) String() string {
	str := ""
	for _, inv := range invs {
		str += fmt.Sprintf("%v\n", *inv)
	}
	return str
}

func InsertOp() string {
	ctx, db := connect()
	_, err := db.Exec(ctx,
		`INSERT INTO invoices (fname, lname, product, price, quantity, category, shipping) `+
			`VALUES ('Larry','Doover', 'Flashlight', 14.99, 5, 'hardware', '543 Kowaoski Road, Salt Lake City, UT 54126')`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	stl := ReadOp()
	return stl
}

// returns the string rslt of  a Select Query
func ReadOp() string {
	ctx, db := connect()
	var invs invoices
	if err := pgxscan.Select(ctx, db, &invs, `SELECT fname, lname, product,
        price, quantity FROM invoices WHERE price > 13.00`); err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	return invs.String()

}

// Create a New Database Connection to bikeshop
func connect() (context.Context, *pgxpool.Pool) {

	uri := "postgres://username@localhost:5432/bikeshop"
	os.Setenv("DATABASE_URL", uri)
	ctx := context.Background()

	db, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to a database: %v\n", err)
		os.Exit(1)
	}

	return ctx, db

}
