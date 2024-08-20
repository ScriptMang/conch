package bikeshop

import (
	"context"
	"fmt"
	"os"
	"strings"

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

// takes an invoice and throws an error for any field with an invalid input
func (inv *Invoice) validateFields() {
	// check for empty fields: for all the fields
	if inv.Fname == "" || inv.Lname == "" || inv.Product == "" ||
		inv.Price == 0.00 || inv.Quantity == 0 || inv.Category == "" || inv.Shipping == "" {
		fmt.Fprintf(os.Stderr, "Error none of the fields can be empty or zero")
		os.Exit(1)
	}

	// check that none of the string fields start or end with a digit or special character
	digitFilter := "0123456789"
	puncFilter := ".,?!'\"`:;"
	symbolFilter := "~@#%$^|><&*()[]{}_-+=\\/"

	// check for digits: first-name, last-name and category
	if strings.IndexAny(inv.Fname, digitFilter) != -1 || strings.IndexAny(inv.Lname, digitFilter) != -1 ||
		strings.IndexAny(inv.Category, digitFilter) != -1 {
		fmt.Fprintf(os.Stderr, "Error the first name, last name or category can't contain any digits")
		os.Exit(1)

	}

	// check for punctuation: first-name, last-name and category
	if strings.IndexAny(inv.Fname, puncFilter) != -1 || strings.IndexAny(inv.Lname, puncFilter) != -1 ||
		strings.IndexAny(inv.Category, puncFilter) != -1 {
		fmt.Fprintf(os.Stderr, "Error the first name, last name or category can't contain any punctuation")
		os.Exit(1)

	}

	// check for spaces: first-name, last-name
	if strings.IndexAny(inv.Fname, " ") != -1 || strings.IndexAny(inv.Lname, " ") != -1 {
		fmt.Fprintf(os.Stderr, "Error the first name, last name or category can't contain any spaces")
		os.Exit(1)

	}

	// check for symbols: first-name, last-name, category
	if strings.IndexAny(inv.Fname, symbolFilter) != -1 || strings.IndexAny(inv.Lname, symbolFilter) != -1 ||
		strings.IndexAny(inv.Category, symbolFilter) != -1 {
		fmt.Fprintf(os.Stderr, "Error the first name, last name or category can't contain any symbols")
		os.Exit(1)

	}

	// check for negative values:  price and quantity
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

// returns all the invoices in the database a slice []*Invoice
func ReadInvoices() []*Invoice {
	ctx, db := connect()
	defer db.Close()

	var invs Invoices
	if err := pgxscan.Select(ctx, db, &invs, `SELECT * FROM invoices`); err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	return invs
}

// return the invoice given the id
// if the id doesn't exist it returns all invoices
func ReadInvoiceByID(id int) Invoice {
	ctx, db := connect()
	defer db.Close()

	row, _ := db.Query(ctx, `SELECT * FROM invoices WHERE id =$1`, id)

	var inv Invoice
	err := pgxscan.ScanOne(&inv, row)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error there is no invoice with the target id of %d: %v\n", id, err)
		os.Exit(1)
	}

	return inv
}

// updates and returns the given invoice by id
func UpdateInvoice(inv Invoice, id int) Invoice {
	ctx, db := connect()
	defer db.Close()

	inv.validateFields()

	qry1 := `UPDATE invoices SET fname=$1,lname=$2,product=$3,` +
		`price=$4,quantity=$5,category=$6,` +
		`shipping=$7 WHERE id=$8 RETURNING *`

	rows, _ := db.Query(ctx, qry1, inv.Fname, inv.Lname, inv.Product,
		inv.Price, inv.Quantity, inv.Category, inv.Shipping, id)

	var inv2 Invoice
	err := pgxscan.ScanOne(&inv2, rows)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error there is no invoice with the target id of %d: %v\n", id, err)
		os.Exit(1)
	}

	return inv2
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

	return ReadInvoices()
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
