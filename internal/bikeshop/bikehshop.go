package bikeshop

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

// meant to be binded to new acct info
type Account struct {
	ID       int    `json:"id" form:"id"`
	Fname    string `json:"fname" form:"fname"`
	Lname    string `json:"lname" form:"lname"`
	Address  string `json:"address" form:"address"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type Users struct {
	ID       int    `json:"id" form:"id"`
	Fname    string `json:"fname" form:"fname"`
	Lname    string `json:"lname" form:"lname"`
	Address  string `json:"address" form:"address"`
	Username string `json:"username" form:"username"`
}

type Passwords struct {
	ID       int    `json:"id" form:"id"`
	UserID   int    `json:"user_id" form:"user_id"`
	Password string `json:"password" form:"password"`
}

type Invoice struct {
	ID       int     `json:"id,omitempty" form:"id,omitempty"`
	UserID   int     `json:"user_id" form:"user_id"`
	Product  string  `json:"product" form:"product"`
	Category string  `json:"category" form:"category"`
	Price    float32 `json:"price" form:"price"`
	Quantity int     `json:"quantity" form:"quantity"`
}

type GrammarError struct {
	ErrMsgs []string
}

// meant to hold err strings related to authentication and account creation
type AuthError struct {
	ErrMsgs []string
}

type Invoices []*Invoice

var ErrorCode int // http-status code for errors
const BadRequest = 400
const resourceNotFound = 404

// helper funct: takes a pointer to an InvoiceErorr, HttpStatusCode and a string msg
// as parameters and sets the values for the GrammarError struct.
// By default content-type is of type 'application/json'
func (fieldErr *GrammarError) AddMsg(statusCode int, str string) {
	ErrorCode = statusCode
	fieldErr.ErrMsgs = append(fieldErr.ErrMsgs, str)
}

// helper funct: takes a pointer to an Authentication Error, HttpStatusCode and a string msg
// as parameters and sets the values for the AuthError struct.
// By default content-type is of type 'application/json'
func (credErr *AuthError) AddMsg(statusCode int, str string) {
	ErrorCode = statusCode
	credErr.ErrMsgs = append(credErr.ErrMsgs, str)
}

// checks for empty text-fields in an invoice
// if there an error its added to an error slice
func isTextFieldEmpty(fieldName string, val *string, fieldErr *GrammarError) {
	if *val == "" {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't be empty")
	}
}

func fieldHasDigits(fieldName string, val *string, fieldErr *GrammarError) {
	digitFilter := "0123456789"
	if isTextInvalid(*val, digitFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any digits")
	}
}

func fieldHasPunct(fieldName string, val *string, fieldErr *GrammarError) {
	punctFilter := ".,?!'\"`:;"

	switch fieldName {
	case "Fname", "Lname":
		punctFilter = " .,?!'\"`:;"
	case "Product":
		punctFilter = "?!'\";"
	case "Category", "Address":
		punctFilter = ".?!'\"`:;"
	}

	if isTextInvalid(*val, punctFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any punctuation")
	}
}

func fieldHasSymbols(fieldName string, val *string, fieldErr *GrammarError) {
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
	if isTextInvalid(*val, symbolFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any Symbols")
	}
}

// checks a string field against an invalid char sequence
// if it returns a index then the text is invalid and it returns true
func isTextInvalid(val, charFilter string) bool {
	return strings.ContainsAny(val, charFilter)
}

// checks the field to see if it exceeds or falls below a given char limit
// if it doesn't match the upper or lower limit an error message is added
// to the list of grammar errors
func isFieldTooLong(fieldName string, val *string, gramErr *GrammarError, minimum, maximum int) {
	fieldLen := len(*val)
	if fieldLen < minimum {
		gramErr.AddMsg(BadRequest, "Error: "+fieldName+" is too short, expected "+
			strconv.Itoa(minimum)+"-"+strconv.Itoa(maximum)+" chars")
	}
	if fieldLen > maximum {
		gramErr.AddMsg(BadRequest, "Error: "+fieldName+" is too long, expected "+
			strconv.Itoa(minimum)+"-"+strconv.Itoa(maximum)+" chars")
	}
}

// checks to see if there any capital letters in string val
// adds an new error to fieldErrs if none exist
func fieldHasNoCapLetters(val *string, fieldErr *GrammarError) {
	capLst := "ABCDEFGHIJKLMNOPQRYTUVWXYZ"
	if !strings.ContainsAny(*val, capLst) {
		fieldErr.AddMsg(BadRequest, "Error: Password must contain one or more capital letters")
	}
}

// checks to see if there are any digits in string val
// adds an new error to fieldErrs if none exist
func fieldHasNoNums(val *string, fieldErr *GrammarError) {
	nums := "0123456789"
	if !strings.ContainsAny(*val, nums) {
		fieldErr.AddMsg(BadRequest, "Error: Password must contain one or more digits")
	}
}

// checks a field for punctuation, digits, and symbols
func checkGrammar(fieldName string, val *string, fieldErr *GrammarError) {

	isTextFieldEmpty(fieldName, val, fieldErr)

	name := fieldName
	if *val != "" && name != "Address" &&
		name != "Product" && name != "Username" &&
		name != "Password" {
		fieldHasDigits(name, val, fieldErr)
		fieldHasPunct(name, val, fieldErr)
		fieldHasSymbols(name, val, fieldErr)
	}

	if name == "Username" ||
		name == "Address" || name == "Product" {
		fieldHasPunct(name, val, fieldErr)
		fieldHasSymbols(name, val, fieldErr)
	}

	if name == "Username" ||
		name == "Password" {
		isFieldTooLong(name, val, fieldErr, 8, 16)
	}

	if name == "Password" {
		fieldHasNoCapLetters(val, fieldErr)
		fieldHasNoNums(val, fieldErr)
	}
	// if name == "Password" {
	// 	fieldHasPunct(field, fieldErr)
	// }
}

// takes an invoice and throws an error for any field with an invalid input
func (inv *Invoice) validateAllFields(user Users) GrammarError {
	// check for empty fields: for all the fields
	textFields := map[string]*string{
		"Fname":    &user.Fname,
		"Lname":    &user.Lname,
		"Category": &inv.Category,
		"Product":  &inv.Product,
		"Address":  &user.Address,
	}
	var fieldErr GrammarError
	for field, val := range textFields {
		checkGrammar(field, val, &fieldErr)
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

// validate username, fname, lname, address fields for digits, symbols, punct
func validateAccount(acct *Account, acctErr *GrammarError) {

	textFields := map[string]*string{
		"Fname":    &acct.Fname,
		"Lname":    &acct.Lname,
		"Address":  &acct.Address,
		"Username": &acct.Username,
		"Password": &acct.Password,
	}

	for field, val := range textFields {
		checkGrammar(field, val, acctErr)
	}

}

func EncryptPassword(val string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(val), 14)
	return string(hash), err
}

// helper funct that adds hash to the passwords table
func AddPassword(acct *Account, acctErr *GrammarError) {
	ctx, db := connect()
	defer db.Close()

	var pswds Passwords
	var err error

	if len(acctErr.ErrMsgs) > 0 {
		return
	}

	// encrypt password
	acct.Password, err = EncryptPassword(acct.Password)
	if err != nil {
		acctErr.AddMsg(BadRequest,
			"Hashing Error: password longer than 72 bytes, can't hash")
		return
	}

	// if no errors add info to appropiate tables
	rows, _ := db.Query(
		ctx,
		`INSERT INTO Passwords (user_id, password) VALUES($1, $2) RETURNING *`,
		acct.ID, acct.Password,
	)

	err = pgxscan.ScanOne(&pswds, rows)
	if err != nil {
		qryError := err.Error()
		if strings.Contains(qryError, "value too long for type character varying") {
			acctErr.AddMsg(BadRequest, "password too long,chars must be less than 72 bytes")
		}
		acctErr.AddMsg(BadRequest, qryError)
	}
}

// helper funct that adds all user-info to users table
func AddUser(acct *Account, acctErr *GrammarError) []*Account {
	ctx, db := connect()
	defer db.Close()

	var insertedAcct Account
	var accts []*Account

	if len(acctErr.ErrMsgs) > 0 {
		fmt.Println("Errs exist in AddUser Funct return nil")
		return nil
	}

	rows, _ := db.Query(
		ctx,
		`INSERT INTO Users (username, fname,lname,address) VALUES($1, $2, $3, $4) RETURNING *`,
		acct.Username, acct.Fname, acct.Lname, acct.Address,
	)

	err := pgxscan.ScanOne(&insertedAcct, rows)
	if err != nil {
		qryError := err.Error()
		if strings.Contains(qryError, "value too long for type character varying") {
			acctErr.AddMsg(BadRequest, "varchar too long, use varchar length between 1-255")
		} else {
			acctErr.AddMsg(BadRequest, qryError)
		}
	}
	//	fmt.Printf("Errors so far when adding a user: %s\n", acctErr.ErrMsgs)
	//	fmt.Printf("New User to be added: %+v\n", insertedAcct)
	accts = append(accts, &insertedAcct)
	return accts

}

// adds the account info to the appropiate tables w/ the database
func AddAccount(acct *Account) ([]*Account, GrammarError) {
	var insertedAcct Account
	var accts []*Account
	acctErr := &GrammarError{}
	validateAccount(acct, acctErr)

	if acctErr.ErrMsgs != nil {
		// fmt.Printf("Errors in AddAccount Func, %v\n", acctErr.ErrMsgs)
		return nil, *acctErr
	}

	// if no errors add info to appropiate tables
	accts = AddUser(acct, acctErr)
	if acctErr.ErrMsgs != nil {
		// fmt.Printf("Errors in AddAccount Func, %v\n", acctErr.ErrMsgs)
		return nil, *acctErr
	}

	// add passwords to table, don't if err existf
	// fmt.Printf("User added into Users: %v\n", *accts[0])
	AddPassword(accts[0], acctErr)
	if acctErr.ErrMsgs != nil {
		// fmt.Printf("Errors in AddAccount Func, %v\n", acctErr.ErrMsgs)
		return nil, *acctErr
	}

	accts = append(accts, &insertedAcct)
	return accts, *acctErr
}

// returns the list of all existing users
func ReadUsers() ([]*Users, GrammarError) {
	ctx, db := connect()
	defer db.Close()

	var usrs []*Users
	fieldErr := GrammarError{ErrMsgs: []string{""}}
	rows, _ := db.Query(ctx, `SELECT * FROM Users`)
	err := pgxscan.ScanAll(&usrs, rows)
	if err != nil {
		errMsg := err.Error()

		if strings.Contains(errMsg, "\"username\" does not exist") {
			fieldErr.ErrMsgs = nil
			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
		}
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}

	return usrs, fieldErr
}

// returns a user given the id
// if the id doesn't exist it error
func ReadUserByID(id int) ([]*Users, GrammarError) {
	ctx, db := connect()
	defer db.Close()

	var usr Users
	var usrs []*Users
	_, fieldErr := ReadUsers()

	// make sure the table isn't empty
	if strings.Contains(fieldErr.ErrMsgs[0], "\"username\" does not exist") {
		return usrs, fieldErr
	}

	row, _ := db.Query(ctx, `SELECT * FROM Users WHERE id=$1`, id)

	err := pgxscan.ScanOne(&usr, row)
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
		return usrs, fieldErr
	}
	usrs = append(usrs, &usr)
	return usrs, fieldErr
}

func InsertOp(usr Users, inv Invoice) ([]*Invoice, GrammarError) {
	ctx, db := connect()
	defer db.Close()

	var insertedInv Invoice
	var invs []*Invoice
	fieldErr := inv.validateAllFields(usr)

	if len(fieldErr.ErrMsgs) > 0 {
		return invs, fieldErr
	}

	// fmt.Printf("The value of the invoice before  InsertOP: %+v\n", inv)
	rows, _ := db.Query(
		ctx,
		`INSERT INTO invoices (user_id, category, product, price, quantity) VALUES($1, $2, $3, $4, $5) RETURNING *`,
		usr.ID, inv.Category, inv.Product, inv.Price, inv.Quantity,
	)

	err := pgxscan.ScanOne(&insertedInv, rows)
	// fmt.Printf("The value of the invoice after InsertOP: %+v\n", &insertedInv)
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

// // returns all the invoices in the database a slice []*Invoice
func ReadInvoices() ([]*Invoice, GrammarError) {
	ctx, db := connect()
	defer db.Close()

	var invs Invoices
	fieldErr := GrammarError{ErrMsgs: []string{""}}
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

func ReadInvoicesByUserID(id int) ([]*Invoice, GrammarError) {
	ctx, db := connect()
	defer db.Close()

	var invs []*Invoice
	_, fieldErr := ReadInvoices()

	if strings.Contains(fieldErr.ErrMsgs[0], "\"username\" does not exist") {
		// fmt.Printf("ReadInvoicesByUserID funct: error username doesn't exist")
		return invs, fieldErr
	}

	rows, _ := db.Query(ctx, `SELECT * FROM invoices WHERE user_id = $1`, id)

	err := pgxscan.ScanAll(&invs, rows)
	if err != nil {
		errMsg := err.Error()
		fieldErr.ErrMsgs = nil
		if strings.Contains(errMsg, "\"username\" does not exist") {
			// fmt.Printf("ReadInvoicesByUserID funct: error username doesn't exist\n")
			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
		} else if strings.Contains(errMsg, "no rows in result set") {
			// fmt.Printf("ReadInvoicesByUserID funct: error invoice with specified id doesn't exist\n")
			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: invoice with specified id does not exist")
		} else {
			// fmt.Printf("ReadInvoicesByUserID funct: error %s\n", err.Error())
			fieldErr.AddMsg(BadRequest, err.Error())

		}

		// fmt.Printf("The len of fieldErr msgs is: %d\n", len(fieldErr.ErrMsgs))
		return invs, fieldErr
	}
	return invs, fieldErr
}

func ReadInvoiceByUserID(userID, invID int) ([]*Invoice, GrammarError) {
	ctx, db := connect()
	defer db.Close()

	var invs []*Invoice
	_, fieldErr := ReadInvoices()

	if strings.Contains(fieldErr.ErrMsgs[0], "\"username\" does not exist") {
		// fmt.Printf("ReadInvoicesByUserID funct: error username doesn't exist")
		return invs, fieldErr
	}

	rows, _ := db.Query(ctx, `SELECT * FROM invoices WHERE user_id = $1 and id = $2`, userID, invID)

	err := pgxscan.ScanAll(&invs, rows)
	if err != nil {
		errMsg := err.Error()
		fieldErr.ErrMsgs = nil
		if strings.Contains(errMsg, "\"username\" does not exist") {
			// fmt.Printf("ReadInvoicesByUserID funct: error username doesn't exist\n")
			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
		} else if strings.Contains(errMsg, "no rows in result set") {
			// fmt.Printf("ReadInvoicesByUserID funct: error invoice with specified id doesn't exist\n")
			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: invoice with specified id does not exist")
		} else {
			// fmt.Printf("ReadInvoicesByUserID funct: error %s\n", err.Error())
			fieldErr.AddMsg(BadRequest, err.Error())

		}

		// fmt.Printf("The len of fieldErr msgs is: %d\n", len(fieldErr.ErrMsgs))
		return invs, fieldErr
	}
	return invs, fieldErr
}

// // return the invoice given the id
// // if the id doesn't exist it returns all invoices
// func ReadInvoiceByID(id int) ([]*Invoice, GrammarError) {
// 	ctx, db := connect()
// 	defer db.Close()

// 	var inv Invoice
// 	var invs []*Invoice
// 	_, fieldErr := ReadInvoices()

// 	if strings.Contains(fieldErr.ErrMsgs[0], "\"username\" does not exist") {
// 		return invs, fieldErr
// 	}

// 	row, _ := db.Query(ctx, `SELECT * FROM invoices WHERE id=$1`, id)

// 	err := pgxscan.ScanOne(&inv, row)
// 	if err != nil {
// 		errMsg := err.Error()
// 		fieldErr.ErrMsgs = nil
// 		if strings.Contains(errMsg, "\"username\" does not exist") {
// 			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
// 		}

// 		if strings.Contains(errMsg, "no rows in result set") {
// 			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: invoice with specified id does not exist")
// 		}

// 		// fmt.Printf("The len of fieldErr msgs is: %d\n", len(fieldErr.ErrMsgs))
// 		return invs, fieldErr
// 	}
// 	invs = append(invs, &inv)
// 	return invs, fieldErr
// }

// func validateFieldsForUpdate(inv *Invoice) GrammarError {
// 	return inv.validateAllFields()
// }

// // updates and returns the given invoice by id
// func UpdateInvoice(inv Invoice, id int) ([]*Invoice, GrammarError) {
// 	ctx, db := connect()
// 	defer db.Close()

// 	var inv2 Invoice // resulting invoice
// 	var invs []*Invoice
// 	var fieldErr GrammarError
// 	_, fieldErr = ReadInvoiceByID(id)
// 	// msgLen := len(fieldErr.ErrMsgs)
// 	// fmt.Printf("There are %d field err messages\n", msgLen)
// 	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
// 		return invs, fieldErr
// 	}

// 	fieldErr = validateFieldsForUpdate(&inv)
// 	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
// 		return invs, fieldErr
// 	}

// 	rows, _ := db.Query(
// 		ctx,
// 		`UPDATE invoices SET fname=$1,lname=$2,product=$3,price=$4,quantity=$5,category=$6,shipping=$7 WHERE id=$8 RETURNING *`,
// 		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping, id,
// 	)

// 	err := pgxscan.ScanOne(&inv2, rows)
// 	if err != nil {
// 		errMsg := err.Error()
// 		fieldErr.ErrMsgs = nil
// 		if strings.Contains(errMsg, "\"username\" does not exist") {
// 			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
// 		} else {
// 			fieldErr.AddMsg(BadRequest, "Invoices are empty")
// 		}
// 		// fmt.Println("%s\n", errMsg)
// 		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
// 	}
// 	invs = append(invs, &inv2)

// 	return invs, fieldErr
// }

// func checkGrammarForPatch(field *textField, orig string, fieldErr *GrammarError) {
// 	name := field.name
// 	if *field.value == "" {
// 		//fmt.Printf("CheckGrammarForPatch: %s field value is blank\n", field.name)
// 		*field.value = orig // unique to patch requests
// 		//fmt.Println("CheckGrammarForPatch: Swap for orig.value: ", field.value)
// 	} else if *field.value != "" && name != "Shipping" && name != "Product" {
// 		fieldHasDigits(*field, fieldErr)
// 		fieldHasPunct(*field, fieldErr)
// 		fieldHasSymbols(*field, fieldErr)
// 	}

// 	if name == "Shipping" || name == "Product" {
// 		fieldHasPunct(*field, fieldErr)
// 		fieldHasSymbols(*field, fieldErr)
// 	}
// }

// func validateFieldsForPatch(orig Invoice, inv *Invoice) GrammarError {
// 	// validate fields for Grammars
// 	modInv := inv
// 	textFields := []*textField{
// 		{name: "Fname", value: &modInv.Fname},
// 		{name: "Lname", value: &modInv.Lname},
// 		{name: "Product", value: &modInv.Product},
// 		{name: "Category", value: &modInv.Category},
// 		{name: "Shipping", value: &modInv.Shipping},
// 	}
// 	var fieldErr GrammarError
// 	origVals := []string{orig.Fname, orig.Lname, orig.Product, orig.Category, orig.Shipping}
// 	for i, text := range textFields {
// 		checkGrammarForPatch(text, origVals[i], &fieldErr)
// 		//fmt.Println("GrammarPatch Returns: ", text.value)
// 		//fmt.Printf("Modified Invoice is: %+v\n", *modInv)
// 	}

// 	if inv.Price == 0 {
// 		inv.Price = orig.Price // unique to patch requests
// 	} else if inv.Price != 0.00 && inv.Price < 0.00 {
// 		fieldErr.AddMsg(BadRequest, "Error: The price can't be negative")
// 		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
// 	}

// 	if inv.Quantity == 0 {
// 		inv.Quantity = orig.Quantity // unique to patch requests
// 	} else if inv.Quantity != 0 && inv.Quantity < 0 {
// 		fieldErr.AddMsg(BadRequest, "Error: The quantity can't be negative")
// 		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
// 	}
// 	return fieldErr
// }

// func PatchInvoice(inv Invoice, id int) ([]*Invoice, GrammarError) {
// 	ctx, db := connect()
// 	defer db.Close()

// 	var inv2 Invoice // resulting invoice
// 	var invs []*Invoice
// 	orig, fieldErr := ReadInvoiceByID(id)
// 	// msgLen := len(fieldErr.ErrMsgs)
// 	// fmt.Printf("There are %d field err messages\n", msgLen)
// 	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
// 		return invs, fieldErr
// 	}

// 	fieldErr = validateFieldsForPatch(*orig[0], &inv)
// 	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
// 		return invs, fieldErr
// 	}

// 	//fmt.Println("PatchInvoice: modified invoice is ", inv)
// 	rows, _ := db.Query(
// 		ctx,
// 		`UPDATE invoices SET fname=$1,lname=$2,product=$3,price=$4,quantity=$5,category=$6,shipping=$7 WHERE id=$8 RETURNING *`,
// 		inv.Fname, inv.Lname, inv.Product, inv.Price, inv.Quantity, inv.Category, inv.Shipping, id,
// 	)

// 	err := pgxscan.ScanOne(&inv2, rows)
// 	if err != nil {
// 		errMsg := err.Error()
// 		fieldErr.ErrMsgs = nil
// 		if strings.Contains(errMsg, "\"username\" does not exist") {
// 			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
// 		} else {
// 			fieldErr.AddMsg(BadRequest, "Invoices are empty")
// 		}
// 		// fmt.Println("%s\n", errMsg)
// 		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
// 	}
// 	invs = append(invs, &inv2)

// 	return invs, fieldErr
// }

// // delete's the given invoice based on id
// // and return the deleted invoice
// func DeleteInvoice(id int) ([]*Invoice, GrammarError) {
// 	ctx, db := connect()
// 	defer db.Close()

// 	var inv Invoice
// 	var invs []*Invoice
// 	_, fieldErr := ReadInvoices()

// 	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
// 		// fmt.Printf("Error messages is empty for Delete-OP")
// 		return invs, fieldErr
// 	}

// 	row, _ := db.Query(ctx, `DELETE FROM invoices WHERE id=$1 RETURNING *`, id)
// 	err := pgxscan.ScanOne(&inv, row)
// 	if err != nil {
// 		errMsg := err.Error()
// 		fieldErr.ErrMsgs = nil
// 		if strings.Contains(errMsg, "\"username\" does not exist") {
// 			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
// 		}

// 		if strings.Contains(errMsg, "no rows in result set") {
// 			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: invoice with specified id does not exist")
// 		}
// 		//fmt.Printf("%s\n", errMsg)
// 		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
// 	}
// 	invs = append(invs, &inv)
// 	return invs, fieldErr
// }

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
