package accts

import (
	"crypto/rand"
	"encoding/hex"
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
	ID      int    `db:"id" json:"id"`
	UserID  int    `db:"user_id" json:"user_id"`
	Fname   string `db:"fname" json:"fname" form:"fname"`
	Lname   string `db:"lname" json:"lname" form:"lname"`
	Address string `db:"address" json:"address"`
}

type Usernames struct {
	ID       int    `db:"id" json:"id"`
	Username string `db:"username" json:"username"`
}

type Passwords struct {
	ID       int    `db:"id" json:"id" form:"id"`
	UserID   int    `db:"user_id" json:"user_id" form:"user_id"`
	Password []byte `db:"password" json:"password" form:"password"`
}

type Tokens struct {
	ID     int    `db:"id" json:"id" form:"id"`
	UserID int    `db:"user_id" json:"user_id"`
	Token  []byte `db:"token" json:"token"`
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

	var newContact UserContacts
	if len(acctErr.ErrMsgs) > 0 {
		// fmt.Println("Errs exist in addUserContact Funct return nil")
		return
	}

	rows, _ := db.Query(
		ctx,
		`INSERT INTO UserContacts (user_id, fname, lname, address) VALUES($1, $2, $3, $4) RETURNING *`,
		acct.ID, acct.Fname, acct.Lname, acct.Address,
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

	var id int
	if len(acctErr.ErrMsgs) > 0 {
		// fmt.Println("Errs exist in AddUser Funct return nil")
		return
	}

	rows, _ := db.Query(
		ctx,
		`INSERT INTO Usernames (username) VALUES($1) RETURNING id`,
		acct.Username,
	)

	err := pgxscan.ScanOne(&id, rows)
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
	acct.ID = id
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

	// // if no errors add info to appropiate tables
	rows, _ := db.Query(ctx,
		`INSERT INTO Passwords (user_id, password) VALUES($1, $2) RETURNING *`,
		acct.ID, hashedPswd,
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
func ReadHashByID(userID int) ([]*Passwords, fields.GrammarError) {
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

// returns random hex as a string
func randHex(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// creates an auth-token for a user and stores it in the database
func GenerateToken(username string, acctErr *fields.GrammarError) Tokens {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var newToken Tokens
	token, _ := randHex(20)

	users, fieldErr := readUserByUsername(username)
	user := *users[0]
	if fieldErr.ErrMsgs != nil {
		return newToken
	}

	rows, _ := db.Query(ctx,
		`INSERT INTO Tokens (user_id, token) VALUES($1, $2) RETURNING *`,
		user.ID, token,
	)

	err := pgxscan.ScanOne(&newToken, rows)
	if err != nil {
		qryError := err.Error()
		switch {
		case strings.Contains(qryError, "value too long for type character varying"):
			acctErr.AddMsg(BadRequest, "varchar too long, use varchar length between 1-255")
		case strings.Contains(qryError, "duplicate key value violates unique constraint"):
			acctErr.AddMsg(BadRequest, "Error: Duplicate User: cannot generate a new token if one already exist")
		default:
			acctErr.AddMsg(BadRequest, qryError)
		}
		return newToken
	}

	return newToken
}

// returns the user id asscoiated by the auth-token
func ReadUserIDByToken(tgtToken string, fieldErr *fields.GrammarError) int {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var token Tokens
	rows, _ := db.Query(ctx,
		`SELECT * FROM tokens WHERE token=$1`, tgtToken,
	)

	err := pgxscan.ScanOne(&token, rows)
	if err != nil {
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "no rows in result set"):
			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: user with specified id does not exist")
		default:
			fieldErr.AddMsg(BadRequest, errMsg)
		}
		return 0
	}
	return token.UserID
}

// adds the account info to the appropiate tables w/ the database
func AddAccount(acct *Account) (*Registered, fields.GrammarError) {
	acctErr := &fields.GrammarError{}
	validateAccount(acct, acctErr)

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

	addUserContact(acct, acctErr)
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

func ReadUsernameByID(userID int, fieldErr *fields.GrammarError) []*Usernames {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var usr Usernames
	var usrs []*Usernames

	row, _ := db.Query(ctx, `SELECT * FROM Usernames WHERE id=$1`, userID)

	err := pgxscan.ScanOne(&usr, row)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "no rows in result set") {
			fieldErr.AddMsg(resourceNotFound, "Resource Not Found: username with specified id doesn't exist")
		}

		// fmt.Printf("The len of fieldErr msgs is: %d\n", len(fieldErr.ErrMsgs))
		return usrs
	}
	usrs = append(usrs, &usr)
	return usrs
}

// When users logout they delete their session token
func LogOut(userID int, fieldErr *fields.GrammarError) Tokens {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var token Tokens
	row, _ := db.Query(ctx,
		`DELETE FROM tokens WHERE user_id=$1 RETURNING *`,
		userID)

	err := pgxscan.ScanOne(&token, row)
	if errors.Is(err, pgx.ErrNoRows) {
		// log.Println("Err: No Rows were Found for the Specified User")
		fieldErr.AddMsg(fields.ResourceNotFound, "Resource Not Found: user with specified id doesn't exist")
		return token
	}

	if err != nil {
		errMsg := err.Error()
		fieldErr.AddMsg(fields.BadRequest, errMsg)
		return token
	}

	return token
}

// Deletes the User account which
// cascades to delete their invoices too
func DeleteAcct(user Usernames) ([]*Usernames, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var usr Usernames
	var usrs []*Usernames

	// verify username exists
	_, fieldErr := readUserByUsername(user.Username)
	if fieldErr.ErrMsgs != nil {
		return nil, fieldErr
	}

	row, _ := db.Query(ctx,
		`DELETE FROM usernames WHERE id=$1 RETURNING *`,
		user.ID)

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
