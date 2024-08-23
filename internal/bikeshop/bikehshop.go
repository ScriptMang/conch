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

type InvoiceError struct {
	ContentType    string
	HttpStatusCode int
	Msg            []string
}

type Invoices []*Invoice

// takes an invoice and throws an error for any field with an invalid input
func (inv *Invoice) validateFields() InvoiceError {
	// check for empty fields: for all the fields
	var fieldErr InvoiceError
	if inv.Fname == "" || inv.Lname == "" || inv.Product == "" ||
		inv.Price == 0.00 || inv.Quantity == 0 || inv.Category == "" || inv.Shipping == "" {
		fieldErr.ContentType = "application/json"
		fieldErr.HttpStatusCode = 400
		fieldErr.Msg = append(fieldErr.Msg, "Error none of the fields can be empty or zero")
		return fieldErr
	}

	// check that none of the string fields start or end with a digit or special character
	digitFilter := "0123456789"
	punctFilter := ".,?!'\"`:;"
	symbolFilter := "~@#%$^|><&*()[]{}_-+=\\/"

	// check for digits: first-name, last-name and category
	if strings.IndexAny(inv.Fname, digitFilter) != -1 || strings.IndexAny(inv.Lname, digitFilter) != -1 ||
		strings.IndexAny(inv.Category, digitFilter) != -1 {
		fieldErr.ContentType = "application/json"
		fieldErr.HttpStatusCode = 400
		fieldErr.Msg = append(fieldErr.Msg, "Bad Request: the first name, last name or category can't contain any digits")
	}

	// check for punctuation: first-name, last-name and category
	if strings.IndexAny(inv.Fname, punctFilter) != -1 || strings.IndexAny(inv.Lname, punctFilter) != -1 ||
		strings.IndexAny(inv.Category, punctFilter) != -1 {
		fieldErr.ContentType = "application/json"
		fieldErr.HttpStatusCode = 400
		fieldErr.Msg = append(fieldErr.Msg, "Bad Request: the first name, last name or category can't contain any punctuation")
	}

	// check for spaces: first-name, last-name
	if strings.IndexAny(inv.Fname, " ") != -1 || strings.IndexAny(inv.Lname, " ") != -1 {
		fieldErr.ContentType = "application/json"
		fieldErr.HttpStatusCode = 400
		fieldErr.Msg = append(fieldErr.Msg, "Bad Request: the first name, last name or category can't contain any spaces")
	}

	// check for symbols: first-name, last-name, category
	if strings.IndexAny(inv.Fname, symbolFilter) != -1 || strings.IndexAny(inv.Lname, symbolFilter) != -1 ||
		strings.IndexAny(inv.Category, symbolFilter) != -1 {
		fieldErr.ContentType = "application/json"
		fieldErr.HttpStatusCode = 400
		fieldErr.Msg = append(fieldErr.Msg, "Bad Request: the first name, last name or category can't contain any symbols")
	}

	// check for negative values:  price and quantity
	if inv.Price < 0.00 || inv.Quantity < 0 {
		fieldErr.ContentType = "application/json"
		fieldErr.HttpStatusCode = 400
		fieldErr.Msg = append(fieldErr.Msg, "Bad Request: Neither the price or quantity can be negative")
	}
	return fieldErr
}

// Takes an invoice and adds it to the database
func InsertOp(inv Invoice) (Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var insertedInv Invoice
	fieldErr := inv.validateFields()

	if len(fieldErr.Msg) > 0 {
		return insertedInv, fieldErr
	}
	rows, _ := db.Query(ctx, `INSERT INTO invoices (fname, lname, product, price, quantity, category, shipping) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *`,
		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping)

	err := pgxscan.ScanOne(&insertedInv, rows)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	return insertedInv, fieldErr
}

// returns all the invoices in the database a slice []*Invoice
func ReadInvoices() []*Invoice {
	ctx, db := connect()
	defer db.Close()

	var invs Invoices
	err := pgxscan.Select(ctx, db, &invs, `SELECT * FROM invoices`)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	return invs
}

// return the invoice given the id
// if the id doesn't exist it returns all invoices
func ReadInvoiceByID(id int) (Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	row, _ := db.Query(ctx, `SELECT * FROM invoices WHERE id=$1`, id)

	var inv Invoice
	fieldErr := inv.validateFields()
	if len(fieldErr.Msg) > 0 {
		return inv, fieldErr
	}

	err := pgxscan.ScanOne(&inv, row)
	if err != nil {
		fieldErr.ContentType = "application/json"
		fieldErr.HttpStatusCode = 404
		fieldErr.Msg = append(fieldErr.Msg, fmt.Sprintf("Resource Not Found: invoice with specified id does not exist: %v\n", err))
	}

	return inv, fieldErr
}

// updates and returns the given invoice by id
func UpdateInvoice(inv Invoice, id int) (Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var inv2 Invoice
	fieldErr := inv.validateFields()
	if len(fieldErr.Msg) > 0 {
		return inv2, fieldErr
	}

	rows, _ := db.Query(ctx, `UPDATE invoices SET fname=$1,lname=$2,product=$3,price=$4,quantity=$5,category=$6,shipping=$7 WHERE id=$8 RETURNING *`,
		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping, id)

	err := pgxscan.ScanOne(&inv2, rows)
	if err != nil {
		fieldErr.ContentType = "application/json"
		fieldErr.HttpStatusCode = 404
		fieldErr.Msg = append(fieldErr.Msg, fmt.Sprintf("Resource Not Found: invoice with specified id does not exist: %v\n", err))
	}

	return inv2, fieldErr
}

// delete's the given invoice based on id
// and return the deleted invoice
func DeleteInvoice(id int) (Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	row, _ := db.Query(ctx, `DELETE FROM invoices WHERE id=$1 RETURNING *`, id)

	var inv Invoice
	fieldErr := inv.validateFields()
	err := pgxscan.ScanOne(&inv, row)
	if err != nil {
		fieldErr.ContentType = "application/json"
		fieldErr.HttpStatusCode = 404
		fieldErr.Msg = append(fieldErr.Msg, fmt.Sprintf("Resource Not Found: invoice with specified id does not exist: %v\n", err))
	}

	return inv, fieldErr
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
