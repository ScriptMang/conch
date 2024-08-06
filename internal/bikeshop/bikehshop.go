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

// Takes an invoice and adds it to the database
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

// returns all the invoices in the database as a slice of *Invoice
func ReadOp() []*Invoice {
	ctx, db := connect()
	var invs invoices
	if err := pgxscan.Select(ctx, db, &invs, `SELECT fname, lname, product,
        price, quantity, category, shipping FROM invoices WHERE price > 13.00`); err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	return invs

}

// updates and returns the given invoice
func UpdateOp(inv Invoice) []*Invoice {
	ctx, db := connect()
	defer db.Close()

	var invs []*Invoice
	set1 := fmt.Sprintf(`SET fname = '%s', lname = '%s', product = '%s', `, inv.Fname, inv.Lname, inv.Product)
	set2 := fmt.Sprintf(`price = %.2f, quantity = %d, category = '%s', `, inv.Price, inv.Quantity, inv.Category)
	set3 := fmt.Sprintf(`shipping = '%s' WHERE fname = '%s'`, inv.Shipping, inv.Fname)
	_, err := db.Exec(ctx, `UPDATE invoices `+set1+set2+set3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}

	qry := fmt.Sprintf(`SELECT * FROM invoices WHERE fname= '%s'`, inv.Fname)
	if err := pgxscan.Select(ctx, db, &invs, qry); err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	return invs
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
