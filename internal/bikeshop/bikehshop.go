package bikeshop

import (
	"context"
	"fmt"
	"os"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Invoice struct {
	ID       int     `json:"id,omitempty" form:"id,omitempty"`
	Fname    string  `json:"fname" form:"fname"`
	Lname    string  `json:"lname" form:"lname"`
	Product  string  `json:"product" form:"product"`
	Price    float32 `json:"price" form:"price"`
	Quantity int     `json:"quantity" form:"quantity"`
	Category string  `json:"category" form:"category"`
	Shipping string  `json:"shipping" form:"shipping"`
}

type Invoices []*Invoice

// prints all the invoices within the slice in json format
func (invs Invoices) Json() string {
	data := ""
	for _, inv := range invs {
		str1 := fmt.Sprintf(`"fname": "%s", "lname": "%s", "product": "%s", `, inv.Fname, inv.Lname, inv.Product)
		str2 := fmt.Sprintf(`"price": %.2f, "quantity": %d, "category": "%s", `, inv.Price, inv.Quantity, inv.Category)
		str3 := fmt.Sprintf(`"shipping": "%s"`, inv.Shipping)
		data += fmt.Sprintf(`{` + str1 + str2 + str3 + `},`)
	}
	return data
}

// prints an invoice in json format
func (inv Invoice) String() string {
	data := ""
	str1 := fmt.Sprintf(`"fname": "%s", "lname": "%s", "product": "%s", `, inv.Fname, inv.Lname, inv.Product)
	str2 := fmt.Sprintf(`"price": %.2f, "quantity": %d, "category": "%s", `, inv.Price, inv.Quantity, inv.Category)
	str3 := fmt.Sprintf(`"shipping": "%s"`, inv.Shipping)
	data += fmt.Sprintf(`{` + str1 + str2 + str3 + `}`)
	return data
}

func (inv *Invoice) validateFields() {
	// make sure none of fields are empty
	if inv.Fname == "" || inv.Lname == "" || inv.Product == "" ||
		inv.Price == 0.00 || inv.Quantity == 0 || inv.Category == "" || inv.Shipping == "" {
		fmt.Fprintf(os.Stderr, "Error none of fields can be empty or zero")
		os.Exit(1)
	}

	//validate none of the numbers are negatives
	if inv.Price < 0.00 || inv.Quantity < 0 {
		fmt.Fprintf(os.Stderr, "Neither the Price or Quantity can be negative")
		os.Exit(1)
	}
}

// Takes an invoice and adds it to the database
func InsertOp(inv Invoice) Invoice {
	ctx, db := connect()
	defer db.Close()

	inv.validateFields()
	rows, _ := db.Query(ctx, `INSERT INTO invoices (fname, lname, product, price, quantity, category, shipping)`+
		`VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *`, inv.Fname, inv.Lname, inv.Product,
		inv.Price, inv.Quantity, inv.Category, inv.Shipping)

	var insertedInv Invoice
	err := pgxscan.ScanOne(&insertedInv, rows)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	return insertedInv
}

// returns all the invoices in the database as a slice of *Invoice
func ReadAllOp() []*Invoice {
	ctx, db := connect()
	defer db.Close()

	var invs Invoices
	if err := pgxscan.Select(ctx, db, &invs, `SELECT id, fname, lname, product,
        price, quantity, category, shipping FROM invoices`); err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	return invs
}

// updates and returns the given invoice
func UpdateOp(inv Invoice) []*Invoice {
	ctx, db := connect()
	defer db.Close()

	inv.validateFields()
	var invs []*Invoice

	qry1 := `UPDATE invoices SET fname=$1,lname=$2,product=$3,` +
		`price=$4,quantity=$5,category=$6,` +
		`shipping=$7 WHERE id=$8`

	_, err := db.Exec(ctx, qry1, inv.Fname, inv.Lname, inv.Product,
		inv.Price, inv.Quantity, inv.Category, inv.Shipping, inv.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}

	qry2 := `SELECT * FROM invoices WHERE id= $1`
	if err := pgxscan.Select(ctx, db, &invs, qry2, inv.ID); err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	return invs
}

// Delete the given invoice based on id
// return the list of remaining invoices
func DeleteOp(inv Invoice) Invoices {
	ctx, db := connect()
	defer db.Close()

	qry := `DELETE FROM invoices WHERE id = $1`
	_, err := db.Exec(ctx, qry, inv.ID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}

	return ReadAllOp()
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
