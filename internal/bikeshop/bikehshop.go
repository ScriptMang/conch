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
	Msg []string
}

type Invoices []*Invoice

var ErrorCode int // http-status code for errors

// helper funct: takes a pointer to an InvoiceErorr, HttpStatusCode and a string msg
// as parameters and sets the values for the InvoiceError struct.
// By default content-type is of type 'application/json'
func (fieldErr *InvoiceError) AddMsg(statusCode int, str string) {
	ErrorCode = statusCode
	fieldErr.Msg = append(fieldErr.Msg, str)
}

// empty-fields
func (fieldErr *InvoiceError) isTextFieldEmpty(field, fieldName string) {
	if field == "" {
		fieldErr.AddMsg(400, "Bad Request, "+fieldName+" field can't be empty")
	}
}

// takes an invoice and throws an error for any field with an invalid input
func (inv *Invoice) validateFields() InvoiceError {
	// check for empty fields: for all the fields
	var fieldErr InvoiceError

	switch {
	case inv.Price == 0.00:
		fieldErr.AddMsg(400, "Bad Request, Price field can't be 0")
	case inv.Quantity == 0:
		fieldErr.AddMsg(400, "Bad Request, Quantity field can't be 0")
	}

	fieldErr.isTextFieldEmpty(inv.Fname, "Fname")
	fieldErr.isTextFieldEmpty(inv.Lname, "Lname")
	fieldErr.isTextFieldEmpty(inv.Product, "Product")
	fieldErr.isTextFieldEmpty(inv.Category, "Category")
	fieldErr.isTextFieldEmpty(inv.Shipping, "Shipping")

	if len(fieldErr.Msg) > 0 {
		return fieldErr
	}

	// check that none of the string fields start or end with a digit or special character
	digitFilter := "0123456789"
	punctFilter := ".,?!'\"`:;"
	symbolFilter := "~@#%$^|><&*()[]{}_-+=\\/"
	productFilter := "?!'\":;~#$|{}_\\+="

	// check for digits: first-name, last-name and category
	if strings.IndexAny(inv.Fname, digitFilter) != -1 || strings.IndexAny(inv.Lname, digitFilter) != -1 ||
		strings.IndexAny(inv.Category, digitFilter) != -1 {
		fieldErr.AddMsg(400, "Bad Request: the first name, last name or category can't contain any digits")
	}

	// check for punctuation: first-name, last-name and category
	if strings.IndexAny(inv.Fname, punctFilter) != -1 || strings.IndexAny(inv.Lname, punctFilter) != -1 ||
		strings.IndexAny(inv.Category, punctFilter) != -1 {
		fieldErr.AddMsg(400, "Bad Request: the first name, last name or category can't contain any punctuation")
	}

	// check specific punctuation and symbols for product
	if strings.IndexAny(inv.Product, productFilter) != -1 {
		if strings.IndexAny(inv.Product, productFilter[:7]) != -1 {
			fieldErr.AddMsg(400, "Bad Request: product can't contain any of the listed forms of punctuation ?!':;\"")
		}
		if strings.IndexAny(inv.Product, productFilter[7:]) != -1 {
			fieldErr.AddMsg(400, "Bad Request: product can't contain any of the listed forms of symbols ~#$|{}_\\+=")
		}
	}

	// check for spaces: first-name, last-name
	if strings.IndexAny(inv.Fname, " ") != -1 || strings.IndexAny(inv.Lname, " ") != -1 {
		fieldErr.AddMsg(400, "Bad Request: the first name, last name or category can't contain any spaces")
	}

	// check for symbols: first-name, last-name, category
	if strings.IndexAny(inv.Fname, symbolFilter) != -1 || strings.IndexAny(inv.Lname, symbolFilter) != -1 ||
		strings.IndexAny(inv.Category, symbolFilter) != -1 {
		fieldErr.AddMsg(400, "Bad Request: the first name, last name or category can't contain any symbols")
	}

	// check for negative values:  price and quantity
	if inv.Price < 0.00 || inv.Quantity < 0 {
		fieldErr.AddMsg(400, "Bad Request: Neither the price or quantity can be negative")
	}
	return fieldErr
}

// Takes an invoice and adds it to the database
func InsertOp(inv Invoice) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var insertedInv Invoice
	var invs []*Invoice
	fieldErr := inv.validateFields()

	if len(fieldErr.Msg) > 0 {
		return invs, fieldErr
	}
	rows, _ := db.Query(ctx, `INSERT INTO invoices (fname, lname, product, price, quantity, category, shipping) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *`,
		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping)

	err := pgxscan.ScanOne(&insertedInv, rows)
	if err != nil {
		qryError := err.Error()
		if strings.Contains(qryError, "numeric field overflow") {
			fieldErr.AddMsg(400, "numeric field overflow, provide a value between 1.00 - 999.99")
		}
		if strings.Contains(qryError, "greater than maximum value for int4") {
			fieldErr.AddMsg(400, "integer overflow, value must be between 1 - 2147483647")
		}
		if strings.Contains(qryError, "value too long for type character varying") {
			fieldErr.AddMsg(400, "varchar too long, use varchar length between 1-255")
		}
		fieldErr.AddMsg(400, qryError)
	}
	invs = append(invs, &insertedInv)

	return invs, fieldErr
}

// returns all the invoices in the database a slice []*Invoice
func ReadInvoices() ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var invs Invoices
	var fieldErr InvoiceError
	err := pgxscan.Select(ctx, db, &invs, `SELECT * FROM invoices`)
	if err != nil {
		fieldErr.AddMsg(400, "Invoices are empty")
	}
	return invs, fieldErr
}

// return the invoice given the id
// if the id doesn't exist it returns all invoices
func ReadInvoiceByID(id int) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	row, _ := db.Query(ctx, `SELECT * FROM invoices WHERE id=$1`, id)

	var inv Invoice
	var invs []*Invoice
	var fieldErr InvoiceError
	err := pgxscan.ScanOne(&inv, row)
	if err != nil {
		fieldErr.AddMsg(404, "Resource Not Found: invoice with specified id does not exist")
		return invs, fieldErr
	}

	invs = append(invs, &inv)
	return invs, fieldErr
}

// updates and returns the given invoice by id
func UpdateInvoice(inv Invoice, id int) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var inv2 Invoice
	var invs []*Invoice
	fieldErr := inv.validateFields()
	if len(fieldErr.Msg) > 0 {
		return invs, fieldErr
	}

	rows, _ := db.Query(ctx, `UPDATE invoices SET fname=$1,lname=$2,product=$3,price=$4,quantity=$5,category=$6,shipping=$7 WHERE id=$8 RETURNING *`,
		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping, id)

	err := pgxscan.ScanOne(&inv2, rows)
	if err != nil {
		fieldErr.AddMsg(404, "Resource Not Found: invoice with specified id does not exist")
	}
	invs = append(invs, &inv2)

	return invs, fieldErr
}

// delete's the given invoice based on id
// and return the deleted invoice
func DeleteInvoice(id int) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	row, _ := db.Query(ctx, `DELETE FROM invoices WHERE id=$1 RETURNING *`, id)

	var inv Invoice
	var invs []*Invoice
	var fieldErr InvoiceError
	err := pgxscan.ScanOne(&inv, row)
	if err != nil {
		fieldErr.AddMsg(404, "Resource Not Found: invoice with specified id does not exist")
	}
	invs = append(invs, &inv)
	return invs, fieldErr
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
