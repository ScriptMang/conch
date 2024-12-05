package accts

import (
	"errors"
	"strings"

	"github.com/ScriptMang/conch/internal/bikeshop"
	"github.com/ScriptMang/conch/internal/fields"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
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

type UserContacts struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Fname   string `json:"fname" form:"fname"`
	Lname   string `json:"lname" form:"lname"`
	Address string `json:"address"`
}

type Usernames struct {
	ID       int    `json:"id" form:"id"`
	Username string `json:"username" form:"username"`
}

type Passwords struct {
	ID       int    `json:"id" form:"id"`
	UserID   int    `json:"user_id" form:"user_id"`
	Password []byte `json:"password" form:"password"`
}

type LoginCred struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"pswd" form:"password"`
}

type Registered struct {
	UserID int
	Msg    string
}
type LoginStatus struct {
	UserID int
	Status string
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

// adds private userinfo  to usercontacts
func addUserContact(acct *Account, acctErr *fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var newContact []*UserContacts
	if len(acctErr.ErrMsgs) > 0 {
		// fmt.Println("Errs exist in addUserContact Funct return nil")
		return
	}

	rows, _ := db.Query(
		ctx,
		`INSERT INTO UserContacts (fname, lname, address) VALUES($1, $2, $3) RETURNING *`,
		acct.Fname, acct.Lname, acct.Address,
	)

	err := pgxscan.ScanOne(&newContact, rows)
	if err != nil {
		qryError := err.Error()
		if strings.Contains(qryError, "value too long for type character varying") {
			acctErr.AddMsg(BadRequest, "varchar too long, use varchar length between 1-255")
		} else {
			acctErr.AddMsg(BadRequest, qryError)
		}
		return
	}
}

// helper funct that adds a users username to users table
func addUsername(acct *Account, acctErr *fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var username []*Usernames

	if len(acctErr.ErrMsgs) > 0 {
		// fmt.Println("Errs exist in AddUser Funct return nil")
		return
	}

	rows, _ := db.Query(
		ctx,
		`INSERT INTO Usernames (username) VALUES($1) RETURNING *`,
		acct.Username,
	)

	err := pgxscan.ScanOne(&username, rows)
	if err != nil {
		qryError := err.Error()
		if strings.Contains(qryError, "value too long for type character varying") {
			acctErr.AddMsg(BadRequest, "varchar too long, use varchar length between 1-255")
		} else {
			acctErr.AddMsg(BadRequest, qryError)
		}
		return
	}

	//	fmt.Printf("Errors so far when adding a user: %s\n", acctErr.ErrMsgs)
	//	fmt.Printf("New User to be added: %+v\n", insertedAcct)
	acct.ID = username[0].ID
}

func encryptPassword(val string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(val), 14)
	return hash, err
}

// helper funct that adds hash to the passwords table
func addPassword(acct *Account, acctErr *fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var pswds Passwords
	// var err error

	if len(acctErr.ErrMsgs) > 0 {
		return
	}

	if len(acct.Password) == 0 {
		acctErr.AddMsg(BadRequest, " Couldn't add password since none exist")
		return
	}
	// encrypt password
	hashedPswd, err := encryptPassword(acct.Password)
	if err != nil {
		acctErr.AddMsg(BadRequest,
			"Hashing Error: password longer than 72 bytes, can't hash")
		return
	}

	// if no errors add info to appropiate tables
	rows, _ := db.Query(
		ctx,
		`INSERT INTO Passwords (password) VALUES($1) RETURNING *`,
		hashedPswd,
	)

	err = pgxscan.ScanOne(&pswds, rows)
	if err != nil {
		qryError := err.Error()
		if strings.Contains(qryError, "value too long for type character varying") {
			acctErr.AddMsg(BadRequest, "password too long,chars must be less than 72 bytes")
		} else {
			acctErr.AddMsg(BadRequest, qryError)
		}
	}
}

// returns a pswd hash given userid
// if the id doesn't exist it error
func readHashByID(userID int) ([]*Passwords, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var pswd Passwords
	var pswds []*Passwords
	_, fieldErr := ReadUserContact()

	// make sure the table isn't empty
	if fieldErr.ErrMsgs != nil {
		return nil, fieldErr
	}

	row, _ := db.Query(ctx, `SELECT Password FROM Passwords WHERE user_id=$1`, userID)

	err := pgxscan.ScanOne(&pswd, row)

	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "no rows in result set"):
			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: user with specified id does not exist")
		default:
			fieldErr.AddMsg(BadRequest, errMsg)
		}
		return nil, fieldErr
	}

	// fmt.Printf("The len of fieldErr msgs is: %d\n", len(fieldErr.ErrMsgs))

	pswds = append(pswds, &pswd)
	return pswds, fieldErr
}

// adds the account info to the appropiate tables w/ the database
func AddAccount(acct *Account) (*Registered, fields.GrammarError) {
	acctErr := &fields.GrammarError{}
	validateAccount(acct, acctErr)

	if acctErr.ErrMsgs != nil {
		// fmt.Printf("Errors in AddAccount Func, %v\n", acctErr.ErrMsgs)
		return nil, *acctErr
	}

	addUserContact(acct, acctErr)
	if acctErr.ErrMsgs != nil {
		// fmt.Printf("Errors in AddAccount Func, %v\n", acctErr.ErrMsgs)
		return nil, *acctErr
	}

	// if no errors add info to appropiate tables
	addUsername(acct, acctErr)
	if acctErr.ErrMsgs != nil {
		// fmt.Printf("Errors in AddAccount Func, %v\n", acctErr.ErrMsgs)
		return nil, *acctErr
	}

	// add passwords to table, don't if err existf
	// fmt.Printf("User added into Usernames: %v\n", *accts[0])
	addPassword(acct, acctErr)
	if acctErr.ErrMsgs != nil {
		// fmt.Printf("Errors in AddAccount Func, %v\n", acctErr.ErrMsgs)
		return nil, *acctErr
	}

	return &Registered{acct.ID, "registered"}, *acctErr
}

// returns the list of all existing users
func ReadUserContact() ([]*UserContacts, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var usrContacts []*UserContacts
	fieldErr := fields.GrammarError{}
	rows, _ := db.Query(ctx, `SELECT * FROM UserContacts`)
	err := pgxscan.ScanAll(&usrContacts, rows)
	if err != nil {
		errMsg := err.Error()
		// log.Printf("Error in ReadUsernames: %s\n", errMsg)
		if strings.Contains(errMsg, "failed to connect") {
			fieldErr.ErrMsgs = nil
			fieldErr.AddMsg(BadRequest, "Error: failed to connect to database, username doesn't exist")
		} else if errMsg != "" {
			fieldErr.AddMsg(BadRequest, errMsg)
		}
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}

	return usrContacts, fieldErr
}

// returns a user given the id
// if the id doesn't exist it error
func ReadUserContactByID(id int) ([]*UserContacts, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var usrContact UserContacts
	var usrContacts []*UserContacts
	_, fieldErr := ReadUserContact()

	// make sure the table isn't empty
	if fieldErr.ErrMsgs != nil {
		return usrContacts, fieldErr
	}

	row, _ := db.Query(ctx, `SELECT * FROM UserContacts WHERE id=$1`, id)

	err := pgxscan.ScanOne(&usrContact, row)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "no rows in result set") {
			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: user with specified id does not exist")
		}

		// fmt.Printf("The len of fieldErr msgs is: %d\n", len(fieldErr.ErrMsgs))
		return usrContacts, fieldErr
	}
	usrContacts = append(usrContacts, &usrContact)
	return usrContacts, fieldErr
}

// returns the user given their username
func readUserByUsername(username string) ([]*Usernames, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var usr Usernames
	var usrs []*Usernames
	_, fieldErr := ReadUserContact()

	// make sure the table isn't empty

	if fieldErr.ErrMsgs != nil {
		return usrs, fieldErr
	}

	row, _ := db.Query(ctx, `SELECT * FROM Usernames WHERE username=$1`, username)

	err := pgxscan.ScanOne(&usr, row)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "no rows in result set") {
			fieldErr.AddMsg(resourceNotFound, "Error: username is incorrect")
		}

		// fmt.Printf("The len of fieldErr msgs is: %d\n", len(fieldErr.ErrMsgs))
		return usrs, fieldErr
	}
	usrs = append(usrs, &usr)
	return usrs, fieldErr
}

// matches the client's username and pswd against the database
// if there's a match the user is logged in, otherwise
// there's an issue with username or password
func LogIntoAcct(userCred LoginCred) (*LoginStatus, fields.GrammarError) {
	// get the user struct check if it exists
	usrs, fieldErr := readUserByUsername(userCred.UserName)
	if fieldErr.ErrMsgs != nil {
		return nil, fieldErr
	}

	// get the stored pswd hash
	usr := usrs[0]
	pswds, fieldErr := readHashByID(usr.ID)
	if fieldErr.ErrMsgs != nil {
		return nil, fieldErr
	}

	// hash the given pswd and compare it to whats
	// stored in the databse
	hashedPswd := pswds[0].Password
	err := bcrypt.CompareHashAndPassword(hashedPswd, []byte(userCred.Password))

	if err != nil {
		// log.Printf("The Password Hash Comparison Failed: %v\n", err.Error())
		fieldErr.AddMsg(BadRequest, "Error: password is incorrect")
		return nil, fieldErr
	}

	return &LoginStatus{usr.ID, "LoggedIn"}, fieldErr
}

// Deletes the User account which cascades
// to delete their invoices too
func DeleteAcct(userCred LoginCred) ([]*Usernames, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var usr Usernames
	var usrs []*Usernames

	// verify that user wants to delete their account by asking
	// for their credentials. if their info matches the db
	// the account will get deleted. Otherwise it returns an
	//error
	isLoggedIn, fieldErr := LogIntoAcct(userCred)
	if fieldErr.ErrMsgs != nil {
		return nil, fieldErr
	}

	row, _ := db.Query(ctx,
		`DELETE FROM users WHERE id=$1 RETURNING *`,
		isLoggedIn.UserID)

	err := pgxscan.ScanOne(&usr, row)
	if errors.Is(err, pgx.ErrNoRows) {
		// log.Println("Err: No Rows were Found for the Specified User")
		fieldErr.AddMsg(fields.ResourceNotFound, "Resource Not Found: user with specified id doesn't exist")
		return nil, fieldErr
	}

	if err != nil {
		// log.Println("Found an Error Iterating in Getting All the Invoices for the Specified User")
		fieldErr.AddMsg(fields.BadRequest, err.Error())
		return nil, fieldErr
	}

	usrs = append(usrs, &usr)
	return usrs, fieldErr

}

// validate username, fname, lname, address fields for digits, symbols, punct
func validateAccount(acct *Account, acctErr *fields.GrammarError) {

	textFields := map[string]*string{
		"Fname":    &acct.Fname,
		"Lname":    &acct.Lname,
		"Address":  &acct.Address,
		"Username": &acct.Username,
		"Password": &acct.Password,
	}

	for field, val := range textFields {
		fields.CheckGrammar(field, val, acctErr)
	}

}
