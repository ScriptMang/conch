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

// checks for empty text-fields in an invoice
// if there an error its added to an error slice
func isTextFieldEmpty(field, fieldName string, fieldErr *InvoiceError) {
	if field == "" {
		fieldErr.AddMsg(400, "Bad Request, "+fieldName+" can't be empty")
	}
}

// check each invoice field for a null value
// if the field is null add error to the invoice error slice
func validateEmptyFields(inv *Invoice, fieldErr *InvoiceError) {

	if inv.Price == 0.00 {
		fieldErr.AddMsg(400, "Bad Request, Price can't be 0")
	}

	if inv.Quantity == 0 {
		fieldErr.AddMsg(400, "Bad Request, Quantity can't be 0")
	}

	isTextFieldEmpty(inv.Fname, "Fname", fieldErr)
	isTextFieldEmpty(inv.Lname, "Lname", fieldErr)
	isTextFieldEmpty(inv.Product, "Product", fieldErr)
	isTextFieldEmpty(inv.Category, "Category", fieldErr)
	isTextFieldEmpty(inv.Shipping, "Shipping", fieldErr)

}

func fieldHasDigits(s, fieldName string, fieldErr *InvoiceError) {
	digitFilter := "0123456789"
	if isTextInvalid(s, digitFilter) {
		fieldErr.AddMsg(400, "Bad Request: "+fieldName+" can't have any digits")
	}
}

func validateFieldsForDigits(inv *Invoice, fieldErr *InvoiceError) {
	// check for digits: first-name, last-name and category

	fieldHasDigits(inv.Fname, "Fname", fieldErr)
	fieldHasDigits(inv.Lname, "Lname", fieldErr)
	fieldHasDigits(inv.Category, "Category", fieldErr)
}

func fieldHasPunct(s string) bool {
	punctFilter := ".,?!'\"`:;"
	return strings.IndexAny(s, punctFilter) != -1
}

// check each invoice field for a punctuation
// if the field has punctuation add error msgs to the invoice error slice msg
func validateFieldsForPunctuation(inv *Invoice, fieldErr *InvoiceError) {
	// check for punctuation: first-name, last-name and category
	if fieldHasPunct(inv.Fname) {
		fieldErr.AddMsg(400, "Bad Request: the first name can't contain any punctuation")
	}

	if fieldHasPunct(inv.Lname) {
		fieldErr.AddMsg(400, "Bad Request: the last name can't contain any punctuation")
	}

	if fieldHasPunct(inv.Category) {
		fieldErr.AddMsg(400, "Bad Request: the category can't contain any punctuation")
	}
}

// checks a string field against an invalid char sequence
// if it returns a index then the text is invalid and it returns true
func isTextInvalid(fieldVal, charFilter string) bool {
	return strings.IndexAny(fieldVal, charFilter) != -1
}

// takes an invoice and throws an error for any field with an invalid input
func (inv *Invoice) validateFields() InvoiceError {
	// check for empty fields: for all the fields
	var fieldErr InvoiceError
	validateEmptyFields(inv, &fieldErr)

	if len(fieldErr.Msg) > 0 {
		return fieldErr
	}

	validateFieldsForDigits(inv, &fieldErr)
	validateFieldsForPunctuation(inv, &fieldErr)

	// check that none of the string fields start or end with a digit or special character

	symbolFilter := "~@#%$^|><&*()[]{}_-+=\\/"
	productFilter := "?!'\":;~#$|{}_\\+="

	// check specific punctuation and symbols for product
	if isTextInvalid(inv.Product, productFilter) {
		if isTextInvalid(inv.Product, productFilter[:7]) {
			fieldErr.AddMsg(400, "Bad Request: product can't contain any of the listed forms of punctuation ?!':;\"")
		}
		if isTextInvalid(inv.Product, productFilter[7:]) {
			fieldErr.AddMsg(400, "Bad Request: product can't contain any of the listed forms of symbols ~#$|{}_\\+=")
		}
	}

	// check for spaces: first-name, last-name
	if isTextInvalid(inv.Fname, " ") || isTextInvalid(inv.Lname, " ") {
		fieldErr.AddMsg(400, "Bad Request: the first name, last name or category can't contain any spaces")
	}

	// check for symbols: first-name, last-name, category
	if isTextInvalid(inv.Fname, symbolFilter) || isTextInvalid(inv.Lname, symbolFilter) ||
		isTextInvalid(inv.Category, symbolFilter) {
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
