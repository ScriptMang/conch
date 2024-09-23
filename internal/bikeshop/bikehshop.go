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
	ErrMsgs []string
}

type Invoices []*Invoice

var ErrorCode int // http-status code for errors
const badRequest = 400

// helper funct: takes a pointer to an InvoiceErorr, HttpStatusCode and a string msg
// as parameters and sets the values for the InvoiceError struct.
// By default content-type is of type 'application/json'
func (fieldErr *InvoiceError) AddMsg(statusCode int, str string) {
	ErrorCode = statusCode
	fieldErr.ErrMsgs = append(fieldErr.ErrMsgs, str)
}

// checks for empty text-fields in an invoice
// if there an error its added to an error slice
func isTextFieldEmpty(field, fieldName string, fieldErr *InvoiceError) {
	if field == "" {
		fieldErr.AddMsg(badRequest, "Bad Request, "+fieldName+" can't be empty")
	}
}

// check each invoice field for a null value
// if the field is null add error to the invoice error slice
func validateAllEmptyFields(inv *Invoice, fieldErr *InvoiceError) {

	if inv.Price == 0.00 {
		fieldErr.AddMsg(badRequest, "Bad Request, Price can't be 0")
	}

	if inv.Quantity == 0 {
		fieldErr.AddMsg(badRequest, "Bad Request, Quantity can't be 0")
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
		fieldErr.AddMsg(badRequest, "Error: "+fieldName+" can't have any digits")
	}
}

// checks fname lname, and category invoice fields for digits
func validateFieldsForDigits(inv *Invoice, fieldErr *InvoiceError) {
	// check for digits: first-name, last-name and category
	fieldHasDigits(inv.Fname, "Fname", fieldErr)
	fieldHasDigits(inv.Lname, "Lname", fieldErr)
	fieldHasDigits(inv.Category, "Category", fieldErr)
}

func fieldHasPunct(s, fieldName string, fieldErr *InvoiceError) {
	punctFilter := ".,?!'\"`:;"

	if fieldName == "Fname" || fieldName == "Lname" {
		punctFilter = " .,?!'\"`:;"
	} else if fieldName == "Product" {
		punctFilter = "?!'\";"
	} else if fieldName == "Category" || fieldName == "Shipping" {
		punctFilter = ".?!'\"`:;"
	}

	if isTextInvalid(s, punctFilter) {
		fieldErr.AddMsg(badRequest, "Error: "+fieldName+" can't have any punctuation")
	}
}

// checks each string invoice field for punctuation
func validateFieldsForPunctuation(inv *Invoice, fieldErr *InvoiceError) {
	fieldHasPunct(inv.Fname, "Fname", fieldErr)
	fieldHasPunct(inv.Lname, "Lname", fieldErr)
	fieldHasPunct(inv.Category, "Category", fieldErr)
	fieldHasPunct(inv.Product, "Product", fieldErr)
	fieldHasPunct(inv.Shipping, "Shipping", fieldErr)
}

func fieldHasSymbols(s, fieldName string, fieldErr *InvoiceError) {
	symbolFilter := "~@#%$^|><&*()[]{}_-+=\\/"

	if fieldName == "Product" {
		symbolFilter = "~#$*{}[]_\\+=><^"
	}

	if fieldName == "Category" {
		symbolFilter = "~@#%$^|><*()[]{}_-+=\\/"
	}

	if fieldName == "Shipping" {
		symbolFilter = "~@#&%$^|><*()[]{}_+=\\/"
	}

	// check for symbols: first-name, last-name, category, product
	if isTextInvalid(s, symbolFilter) {
		fieldErr.AddMsg(badRequest, "Error: "+fieldName+" can't have any Symbols")
	}
}

// checks each string invoice field for symbols
func validateFieldsForSymbols(inv *Invoice, fieldErr *InvoiceError) {
	fieldHasSymbols(inv.Fname, "Fname", fieldErr)
	fieldHasSymbols(inv.Lname, "Lname", fieldErr)
	fieldHasSymbols(inv.Category, "Category", fieldErr)
	fieldHasSymbols(inv.Product, "Product", fieldErr)
	fieldHasSymbols(inv.Shipping, "Shipping", fieldErr)
}

// checks a string field against an invalid char sequence
// if it returns a index then the text is invalid and it returns true
func isTextInvalid(fieldVal, charFilter string) bool {
	return strings.IndexAny(fieldVal, charFilter) != -1
}

// takes an invoice and throws an error for any field with an invalid input
func (inv *Invoice) validateAllFields() InvoiceError {
	// check for empty fields: for all the fields
	var fieldErr InvoiceError
	validateAllEmptyFields(inv, &fieldErr)

	if len(fieldErr.ErrMsgs) > 0 {
		return fieldErr
	}

	validateFieldsForDigits(inv, &fieldErr)
	validateFieldsForPunctuation(inv, &fieldErr)
	validateFieldsForSymbols(inv, &fieldErr)

	// check for negative values:  price and quantity
	if inv.Price < 0.00 || inv.Quantity < 0 {
		fieldErr.AddMsg(badRequest, "Error: Neither the price or quantity can be negative")
	}
	return fieldErr
}

// Takes an invoice and adds it to the database
func InsertOp(inv Invoice) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var insertedInv Invoice
	var invs []*Invoice
	fieldErr := inv.validateAllFields()

	if len(fieldErr.ErrMsgs) > 0 {
		return invs, fieldErr
	}
	rows, _ := db.Query(ctx, `INSERT INTO invoices (fname, lname, product, price, quantity, category, shipping) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *`,
		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping)

	err := pgxscan.ScanOne(&insertedInv, rows)
	if err != nil {
		qryError := err.Error()
		if strings.Contains(qryError, "numeric field overflow") {
			fieldErr.AddMsg(badRequest, "numeric field overflow, provide a value between 1.00 - 999.99")
		}
		if strings.Contains(qryError, "greater than maximum value for int4") {
			fieldErr.AddMsg(badRequest, "integer overflow, value must be between 1 - 2147483647")
		}
		if strings.Contains(qryError, "value too long for type character varying") {
			fieldErr.AddMsg(badRequest, "varchar too long, use varchar length between 1-255")
		}
		fieldErr.AddMsg(badRequest, qryError)
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
	//rows, _ := db.Query(ctx, `SELECT * FROM invoices`)
	rows, _ := db.Query(ctx, `SELECT * FROM invoices`)
	err := pgxscan.ScanAll(&invs, rows)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.AddMsg(badRequest, "Error: failed to connect to database, username doesn't exist")
		}
	}

	if invs == nil && err == nil {
		fieldErr.AddMsg(400, "Error: The database table invoices is empty")
	}

	return invs, fieldErr
}

// return the invoice given the id
// if the id doesn't exist it returns all invoices
func ReadInvoiceByID(id int) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var inv Invoice
	var invs []*Invoice
	_, fieldErr := ReadInvoices()
	if len(fieldErr.ErrMsgs) > 0 {
		return invs, fieldErr
	}

	row, _ := db.Query(ctx, `SELECT * FROM invoices WHERE id=$1`, id)

	err := pgxscan.ScanOne(&inv, row)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.AddMsg(badRequest, "Error: failed to connect to database, username doesn't exist")
		}

		if strings.Contains(errMsg, "no rows in result set") {
			fieldErr.AddMsg(404, "Resource Not Found: invoice with specified id does not exist")
		}
		fmt.Printf("%s\n", errMsg)
		return invs, fieldErr
	}

	invs = append(invs, &inv)
	return invs, fieldErr
}

func checkGrammar(val *string, orig, fieldName string, fieldErr *InvoiceError) {
	if *val == "" {
		*val = orig
	} else if *val != "" && fieldName != "Shipping" && fieldName != "Product" {
		fieldHasDigits(*val, fieldName, fieldErr)
		fieldHasPunct(*val, fieldName, fieldErr)
		fieldHasSymbols(*val, fieldName, fieldErr)
	}

	if fieldName == "Shipping" || fieldName == "Product" {
		fieldHasPunct(*val, fieldName, fieldErr)
		fieldHasSymbols(*val, fieldName, fieldErr)
	}
}

func validateFieldsForUpdate(orig Invoice, inv *Invoice) InvoiceError {
	var fieldErr InvoiceError

	//validate fields for Grammars, ignore if empty

	checkGrammar(&inv.Fname, orig.Fname, "Fname", &fieldErr)
	checkGrammar(&inv.Lname, orig.Lname, "Lname", &fieldErr)
	checkGrammar(&inv.Product, orig.Product, "Product", &fieldErr)
	checkGrammar(&inv.Category, orig.Category, "Category", &fieldErr)
	checkGrammar(&inv.Shipping, orig.Shipping, "Shipping", &fieldErr)

	if inv.Price == 0 {
		inv.Price = orig.Price
	} else if inv.Price != 0.00 && inv.Price < 0.00 {
		fieldErr.AddMsg(badRequest, "Error: The price can't be negative")
	}

	if inv.Quantity == 0 {
		inv.Quantity = orig.Quantity
	} else if inv.Quantity != 0 && inv.Quantity < 0 {
		fieldErr.AddMsg(badRequest, "Error: The quantity can't be negative")
	}
	return fieldErr
}

// updates and returns the given invoice by id
func UpdateInvoice(inv Invoice, id int) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var orig []*Invoice // original invoice
	var inv2 Invoice    // resulting invoice
	var invs []*Invoice
	var fieldErr InvoiceError
	orig, fieldErr = ReadInvoiceByID(id)
	msgLen := len(fieldErr.ErrMsgs)
	fmt.Printf("There are %d field err messages\n", msgLen)
	if len(fieldErr.ErrMsgs) > 0 {
		return invs, fieldErr
	}

	fieldErr = validateFieldsForUpdate(*orig[0], &inv)
	if len(fieldErr.ErrMsgs) > 0 {
		return invs, fieldErr
	}

	rows, _ := db.Query(ctx, `UPDATE invoices SET fname=$1,lname=$2,product=$3,price=$4,quantity=$5,category=$6,shipping=$7 WHERE id=$8 RETURNING *`,
		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping, id)

	err := pgxscan.ScanOne(&inv2, rows)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.AddMsg(badRequest, "Error: failed to connect to database, username doesn't exist")
		} else {
			fieldErr.AddMsg(badRequest, "Invoices are empty")
		}
		fmt.Println("%s\n", errMsg)
	}
	invs = append(invs, &inv2)

	return invs, fieldErr
}

// delete's the given invoice based on id
// and return the deleted invoice
func DeleteInvoice(id int) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var inv Invoice
	var invs []*Invoice
	_, fieldErr := ReadInvoices()
	if len(fieldErr.ErrMsgs) > 0 {
		return invs, fieldErr
	}

	row, _ := db.Query(ctx, `DELETE FROM invoices WHERE id=$1 RETURNING *`, id)
	err := pgxscan.ScanOne(&inv, row)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.AddMsg(badRequest, "Error: failed to connect to database, username doesn't exist")
		}

		if strings.Contains(errMsg, "no rows in result set") {
			fieldErr.AddMsg(404, "Resource Not Found: invoice with specified id does not exist")
		}
		fmt.Printf("%s\n", errMsg)
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
