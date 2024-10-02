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
type textField struct {
	name  string // field-name
	value string // field-value
}
type InvoiceError struct {
	ErrMsgs []string
}

type Invoices []*Invoice

var ErrorCode int // http-status code for errors
const BadRequest = 400
const resourceNotFound = 404

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
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't be empty")
	}
}

// check each invoice field for a null value
// if the field is null add error to the invoice error slice
func validateAllEmptyFields(inv *Invoice, fieldErr *InvoiceError) {

	if inv.Price == 0.00 {
		fieldErr.AddMsg(BadRequest, "Error: Price can't be zero")
	}

	if inv.Quantity == 0 {
		fieldErr.AddMsg(BadRequest, "Error: Quantity can't be zero")
	}

	textFields := []textField{{"Fname", inv.Fname}, {"Lname", inv.Lname},
		{"Category", inv.Category}, {"Product", inv.Product}, {"Shipping", inv.Shipping}}
	for _, text := range textFields {
		isTextFieldEmpty(text.value, text.name, fieldErr)
	}
}

func fieldHasDigits(s, fieldName string, fieldErr *InvoiceError) {
	digitFilter := "0123456789"
	if isTextInvalid(s, digitFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any digits")
	}
}

// checks fname lname, and category invoice fields for digits
func validateFieldsForDigits(inv *Invoice, fieldErr *InvoiceError) {
	// check for digits: first-name, last-name and category
	textFields := []textField{{"Fname", inv.Fname}, {"Lname", inv.Lname}, {"Category", inv.Category}}
	for _, text := range textFields {
		fieldHasDigits(text.value, text.name, fieldErr)
	}
}

func fieldHasPunct(s, fieldName string, fieldErr *InvoiceError) {
	punctFilter := ".,?!'\"`:;"

	switch fieldName {
	case "Fname", "Lname":
		punctFilter = " .,?!'\"`:;"
	case "Product":
		punctFilter = "?!'\";"
	case "Category", "Shipping":
		punctFilter = ".?!'\"`:;"
	}

	if isTextInvalid(s, punctFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any punctuation")
	}
}

// checks each string invoice field for punctuation
func validateFieldsForPunctuation(inv *Invoice, fieldErr *InvoiceError) {
	textFields := []textField{{"Fname", inv.Fname}, {"Lname", inv.Lname},
		{"Category", inv.Category}, {"Product", inv.Product}, {"Shipping", inv.Shipping}}
	for _, text := range textFields {
		fieldHasPunct(text.value, text.name, fieldErr)
	}
}

func fieldHasSymbols(s, fieldName string, fieldErr *InvoiceError) {
	symbolFilter := "~@#%$^|><&*()[]{}_-+=\\/"

	switch fieldName {
	case "Product":
		symbolFilter = "~#$*{}[]_\\+=><^"
	case "Category":
		symbolFilter = "~@#%$^|><*()[]{}_-+=\\/"
	case "Shipping":
		symbolFilter = "~@#&%$^|><*()[]{}_+=\\/"
	}

	// check for symbols: first-name, last-name, category, product
	if isTextInvalid(s, symbolFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any Symbols")
	}
}

// // checks each string invoice field for symbols
func validateFieldsForSymbols(inv *Invoice, fieldErr *InvoiceError) {
	textFields := []textField{{"Fname", inv.Fname}, {"Lname", inv.Lname},
		{"Category", inv.Category}, {"Product", inv.Product}, {"Shipping", inv.Shipping}}
	for _, text := range textFields {
		fieldHasSymbols(text.value, text.name, fieldErr)
	}
}

// checks a string field against an invalid char sequence
// if it returns a index then the text is invalid and it returns true
func isTextInvalid(fieldVal, charFilter string) bool {
	return strings.IndexAny(fieldVal, charFilter) != -1
}

// checks a field for punctuation, digits, and symbols
func checkGrammar(val *string, fieldName string, fieldErr *InvoiceError) {

	isTextFieldEmpty(*val, fieldName, fieldErr)

	if *val != "" && fieldName != "Shipping" && fieldName != "Product" {
		fieldHasDigits(*val, fieldName, fieldErr)
		fieldHasPunct(*val, fieldName, fieldErr)
		fieldHasSymbols(*val, fieldName, fieldErr)
	}

	if fieldName == "Shipping" || fieldName == "Product" {
		fieldHasPunct(*val, fieldName, fieldErr)
		fieldHasSymbols(*val, fieldName, fieldErr)
	}
}

// takes an invoice and throws an error for any field with an invalid input
func (inv *Invoice) validateAllFields() InvoiceError {
	// check for empty fields: for all the fields
	textFields := []textField{
		{name: "Fname", value: inv.Fname},
		{name: "Lname", value: inv.Lname},
		{name: "Category", value: inv.Category},
		{name: "Product", value: inv.Product},
		{name: "Shipping", value: inv.Shipping},
	}
	var fieldErr InvoiceError
	for _, text := range textFields {
		checkGrammar(&text.value, text.name, &fieldErr)
	}

	// check for negative values:  price and quantity
	if inv.Price == 0.00 {
		fieldErr.AddMsg(BadRequest, "Error: Price can't be zero")
	} else if inv.Price < 0.00 {
		fieldErr.AddMsg(BadRequest, "Error: The price can't be negative")
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}

	if inv.Quantity == 0 {
		fieldErr.AddMsg(BadRequest, "Error: Quantity can't be zero")
	} else if inv.Quantity < 0 {
		fieldErr.AddMsg(BadRequest, "Error: The quantity can't be negative")
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
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

	rows, _ := db.Query(
		ctx,
		`INSERT INTO invoices (fname, lname, product, price, quantity, category, shipping) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING *`,
		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping,
	)

	err := pgxscan.ScanOne(&insertedInv, rows)
	if err != nil {
		qryError := err.Error()
		if strings.Contains(qryError, "numeric field overflow") {
			fieldErr.AddMsg(BadRequest, "numeric field overflow, provide a value between 1.00 - 999.99")
		}
		if strings.Contains(qryError, "greater than maximum value for int4") {
			fieldErr.AddMsg(BadRequest, "integer overflow, value must be between 1 - 2147483647")
		}
		if strings.Contains(qryError, "value too long for type character varying") {
			fieldErr.AddMsg(BadRequest, "varchar too long, use varchar length between 1-255")
		}
		fieldErr.AddMsg(BadRequest, qryError)
	}
	invs = append(invs, &insertedInv)

	return invs, fieldErr
}

// returns all the invoices in the database a slice []*Invoice
func ReadInvoices() ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var invs Invoices
	fieldErr := InvoiceError{ErrMsgs: []string{""}}
	rows, _ := db.Query(ctx, `SELECT * FROM invoices`)
	err := pgxscan.ScanAll(&invs, rows)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.ErrMsgs = nil
			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
		}
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
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

	if strings.Contains(fieldErr.ErrMsgs[0], "\"username\" does not exist") {
		return invs, fieldErr
	}

	row, _ := db.Query(ctx, `SELECT * FROM invoices WHERE id=$1`, id)

	err := pgxscan.ScanOne(&inv, row)
	if err != nil {
		errMsg := err.Error()
		fieldErr.ErrMsgs = nil
		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
		}

		if strings.Contains(errMsg, "no rows in result set") {
			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: invoice with specified id does not exist")
		}

		// fmt.Printf("The len of fieldErr msgs is: %d\n", len(fieldErr.ErrMsgs))
		return invs, fieldErr
	}
	invs = append(invs, &inv)
	return invs, fieldErr
}

func validateFieldsForUpdate(inv *Invoice) InvoiceError {
	var fieldErr InvoiceError

	// validate fields for Grammars, ignore if empty

	checkGrammar(&inv.Fname, "Fname", &fieldErr)
	checkGrammar(&inv.Lname, "Lname", &fieldErr)
	checkGrammar(&inv.Product, "Product", &fieldErr)
	checkGrammar(&inv.Category, "Category", &fieldErr)
	checkGrammar(&inv.Shipping, "Shipping", &fieldErr)

	if inv.Price == 0.00 {
		fieldErr.AddMsg(BadRequest, "Error: Price can't be zero")
	} else if inv.Price < 0.00 {
		fieldErr.AddMsg(BadRequest, "Error: The price can't be negative")
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}

	if inv.Quantity == 0 {
		fieldErr.AddMsg(BadRequest, "Error: Quantity can't be zero")
	} else if inv.Quantity < 0 {
		fieldErr.AddMsg(BadRequest, "Error: The quantity can't be negative")
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}
	return fieldErr
}

// updates and returns the given invoice by id
func UpdateInvoice(inv Invoice, id int) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var inv2 Invoice // resulting invoice
	var invs []*Invoice
	var fieldErr InvoiceError
	_, fieldErr = ReadInvoiceByID(id)
	// msgLen := len(fieldErr.ErrMsgs)
	// fmt.Printf("There are %d field err messages\n", msgLen)
	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		return invs, fieldErr
	}

	fieldErr = validateFieldsForUpdate(&inv)
	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		return invs, fieldErr
	}

	rows, _ := db.Query(
		ctx,
		`UPDATE invoices SET fname=$1,lname=$2,product=$3,price=$4,quantity=$5,category=$6,shipping=$7 WHERE id=$8 RETURNING *`,
		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping, id,
	)

	err := pgxscan.ScanOne(&inv2, rows)
	if err != nil {
		errMsg := err.Error()
		fieldErr.ErrMsgs = nil
		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
		} else {
			fieldErr.AddMsg(BadRequest, "Invoices are empty")
		}
		// fmt.Println("%s\n", errMsg)
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}
	invs = append(invs, &inv2)

	return invs, fieldErr
}

func checkGrammarForPatch(val *string, orig, fieldName string, fieldErr *InvoiceError) {
	if *val == "" {
		*val = orig // unique to patch requests
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

func validateFieldsForPatch(orig Invoice, inv *Invoice) InvoiceError {
	var fieldErr InvoiceError

	// validate fields for Grammars, ignore if empty

	checkGrammarForPatch(&inv.Fname, orig.Fname, "Fname", &fieldErr)
	checkGrammarForPatch(&inv.Lname, orig.Lname, "Lname", &fieldErr)
	checkGrammarForPatch(&inv.Product, orig.Product, "Product", &fieldErr)
	checkGrammarForPatch(&inv.Category, orig.Category, "Category", &fieldErr)
	checkGrammarForPatch(&inv.Shipping, orig.Shipping, "Shipping", &fieldErr)

	if inv.Price == 0 {
		inv.Price = orig.Price // unique to patch requests
	} else if inv.Price != 0.00 && inv.Price < 0.00 {
		fieldErr.AddMsg(BadRequest, "Error: The price can't be negative")
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}

	if inv.Quantity == 0 {
		inv.Quantity = orig.Quantity // unique to patch requests
	} else if inv.Quantity != 0 && inv.Quantity < 0 {
		fieldErr.AddMsg(BadRequest, "Error: The quantity can't be negative")
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}
	return fieldErr
}

func PatchInvoice(inv Invoice, id int) ([]*Invoice, InvoiceError) {
	ctx, db := connect()
	defer db.Close()

	var inv2 Invoice // resulting invoice
	var invs []*Invoice
	orig, fieldErr := ReadInvoiceByID(id)
	// msgLen := len(fieldErr.ErrMsgs)
	// fmt.Printf("There are %d field err messages\n", msgLen)
	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		return invs, fieldErr
	}

	fieldErr = validateFieldsForPatch(*orig[0], &inv)
	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		return invs, fieldErr
	}

	rows, _ := db.Query(
		ctx,
		`UPDATE invoices SET fname=$1,lname=$2,product=$3,price=$4,quantity=$5,category=$6,shipping=$7 WHERE id=$8 RETURNING *`,
		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping, id,
	)

	err := pgxscan.ScanOne(&inv2, rows)
	if err != nil {
		errMsg := err.Error()
		fieldErr.ErrMsgs = nil
		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
		} else {
			fieldErr.AddMsg(BadRequest, "Invoices are empty")
		}
		// fmt.Println("%s\n", errMsg)
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
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

	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		// fmt.Printf("Error messages is empty for Delete-OP")
		return invs, fieldErr
	}

	row, _ := db.Query(ctx, `DELETE FROM invoices WHERE id=$1 RETURNING *`, id)
	err := pgxscan.ScanOne(&inv, row)
	if err != nil {
		errMsg := err.Error()
		fieldErr.ErrMsgs = nil
		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
		}

		if strings.Contains(errMsg, "no rows in result set") {
			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: invoice with specified id does not exist")
		}
		//fmt.Printf("%s\n", errMsg)
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
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
