package bikeshop

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ScriptMang/conch/internal/invs"
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

// meant to hold err strings related to authentication and account creation
type AuthError struct {
	ErrMsgs []string
}

var ErrorCode int // http-status code for errors
const BadRequest = 400
const resourceNotFound = 404

// helper funct: takes a pointer to an Authentication Error, HttpStatusCode and a string msg
// as parameters and sets the values for the AuthError struct.
// By default content-type is of type 'application/json'
func (credErr *AuthError) AddMsg(statusCode int, str string) {
	ErrorCode = statusCode
	credErr.ErrMsgs = append(credErr.ErrMsgs, str)
}

// checks for empty text-fields in an invoice
// if there an error its added to an error slice
func isTextFieldEmpty(fieldName string, val *string, fieldErr *invs.GrammarError) {
	if *val == "" {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't be empty")
	}
}

func fieldHasDigits(fieldName string, val *string, fieldErr *invs.GrammarError) {
	digitFilter := "0123456789"
	if isTextInvalid(*val, digitFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any digits")
	}
}

func fieldHasPunct(fieldName string, val *string, fieldErr *invs.GrammarError) {
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

func fieldHasSymbols(fieldName string, val *string, fieldErr *invs.GrammarError) {
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
func isFieldTooLong(fieldName string, val *string, gramErr *invs.GrammarError, minimum, maximum int) {
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
func fieldHasNoCapLetters(val *string, fieldErr *invs.GrammarError) {
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
func checkGrammar(fieldName string, val *string, fieldErr *invs.GrammarError) {

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

// validate username, fname, lname, address fields for digits, symbols, punct
func validateAccount(acct *Account, acctErr *invs.GrammarError) {

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
func AddUser(acct *Account, acctErr *invs.GrammarError) []*Account {
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
func AddAccount(acct *Account) ([]*Account, invs.GrammarError) {
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
func ReadUsers() ([]*Users, invs.GrammarError) {
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
func ReadUserByID(id int) ([]*Users, invs.GrammarError) {
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

// Create a New Database Connection to bikeshop
func Connect() (context.Context, *pgxpool.Pool) {
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
