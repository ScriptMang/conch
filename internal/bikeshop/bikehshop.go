package bikeshop

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Passwords struct {
	ID       int    `json:"id" form:"id"`
	UserID   int    `json:"user_id" form:"user_id"`
	Password string `json:"password" form:"password"`
}

var ErrorCode int // http-status code for errors
const BadRequest = 400
const resourceNotFound = 404

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
