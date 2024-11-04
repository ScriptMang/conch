package accts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
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
