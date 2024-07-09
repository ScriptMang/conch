package bikeshop

import (
	"context"
	"fmt"
	"os"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Invoice struct {
	Fname    string  `json:"fname" form:"fname"`
	Lname    string  `json:"lname" form:"lname"`
	Product  string  `json:"product" form:"product"`
	Price    float32 `json:"price" form:"price"`
	Quantity int     `json:"quantity" form:"quantity"`
	Category string  `json:"category" form:"category"`
	Shipping string  `json:"shipping" form:"shipping"`
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

// Receive the new tableRow from Post-Rqst and execute the insert-op
func InsertOp(inv Invoice) {
	ctx, db := connect()

	columns := fmt.Sprint(`(fname, lname, product, price, quantity, category, shipping)`)
	rowData := fmt.Sprintf(`('%s','%s','%s',%.2f,%d,'%s','%s')`,
		inv.Fname, inv.Lname, inv.Product, inv.Price,
		inv.Quantity, inv.Category, inv.Shipping)
	_, err := db.Exec(ctx, `INSERT INTO invoices `+columns+`VALUES `+rowData)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
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
