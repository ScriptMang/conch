package bikeshop

import (
	"context"
	"fmt"
	"os"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Invoice struct {
	Fname    string  `json:"fname" form:"fname"`
	Lname    string  `json:"lname" form:"lname"`
	Product  string  `json:"product" form:"product"`
	Price    float32 `json:"price" form:"price"`
	Quantity int     `json:"quantity" form:"quantity"`
	Category string  `json:"category" form:"category"`
	Shipping string  `json:"shipping" form:"shipping"`
}

// invoice meant for deletion
type DeletionForm struct {
	Val string `form:"invoice_list"`
}

type Invoices []*Invoice

// prints all the invoices within the slice in json format
func (invs Invoices) Json() string {
	data := ""
	for _, inv := range invs {
		str1 := fmt.Sprintf(`"fname": "%s", "lname": "%s", "product": "%s", `, inv.Fname, inv.Lname, inv.Product)
		str2 := fmt.Sprintf(`"price": %.2f, "quantity": %d, "category": "%s", `, inv.Price, inv.Quantity, inv.Category)
		str3 := fmt.Sprintf(`"shipping": "%s"`, inv.Shipping)
		data += fmt.Sprintf(`{` + str1 + str2 + str3 + `},`)
	}
	return data
}

// prints an invoice in json format
func (inv Invoice) String() string {
	data := ""
	str1 := fmt.Sprintf(`"fname": "%s", "lname": "%s", "product": "%s", `, inv.Fname, inv.Lname, inv.Product)
	str2 := fmt.Sprintf(`"price": %.2f, "quantity": %d, "category": "%s", `, inv.Price, inv.Quantity, inv.Category)
	str3 := fmt.Sprintf(`"shipping": "%s"`, inv.Shipping)
	data += fmt.Sprintf(`{` + str1 + str2 + str3 + `}`)
	return data
}

// adds an invoice's properties to a sql query if the values given
// from the htlm-form are not blank
func addAttribsToQuery(qry *string, attribName, val string) {
	if val != "" {
		*qry += fmt.Sprintf(`%s = '%s',`, attribName, val)
	}
}

// Takes an invoice and adds it to the database
func InsertOp(inv Invoice) {
	ctx, db := connect()

	columns := fmt.Sprint(`(fname, lname, product, price, quantity, category, shipping)`)
	rowData := fmt.Sprintf(`('%s','%s','%s',%.2f,%d,'%s','%s')`,
		inv.Fname, inv.Lname, inv.Product, inv.Price,
		inv.Quantity, inv.Category, inv.Shipping)
	_, err := db.Exec(ctx, `INSERT INTO invoices `+columns+`VALUES `+rowData)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
}

// returns all the invoices in the database as a slice of *Invoice
func ReadOp() []*Invoice {
	ctx, db := connect()
	var invs Invoices
	if err := pgxscan.Select(ctx, db, &invs, `SELECT fname, lname, product,
        price, quantity, category, shipping FROM invoices`); err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	return invs

}

// updates and returns the given invoice
func UpdateOp(inv Invoice) []*Invoice {
	ctx, db := connect()
	defer db.Close()

	var invs []*Invoice

	cols := ""

	addAttribsToQuery(&cols, `fname`, inv.Fname)
	addAttribsToQuery(&cols, `lname`, inv.Lname)
	addAttribsToQuery(&cols, `product`, inv.Product)

	if inv.Price > 0.00 {
		cols += fmt.Sprintf(`price = %.2f,`, inv.Price)
	}

	if inv.Quantity > 0 {
		cols += fmt.Sprintf(`quantity = %d,`, inv.Quantity)
	}

	addAttribsToQuery(&cols, `category`, inv.Category)
	addAttribsToQuery(&cols, `shipping`, inv.Shipping)

	// remove the apostrophe at the end of the columns
	stripCols := cols
	lastCharIdx := len(cols) - 1
	if lastCharIdx >= 0 && cols[lastCharIdx] == ',' {
		stripCols = cols[:lastCharIdx]
	}

	expr := `SET ` + stripCols + fmt.Sprintf(` WHERE fname = '%s'`, inv.Fname)
	fmt.Println("This is the update query w/o set stripped: ", expr)

	_, err := db.Exec(ctx, `UPDATE invoices `+expr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}

	qry := fmt.Sprintf(`SELECT * FROM invoices WHERE fname= '%s'`, inv.Fname)
	if err := pgxscan.Select(ctx, db, &invs, qry); err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}
	return invs
}

// Delete the given invoice based on fname
// return the list of remaining invoices
func DeleteOp(fname string) Invoices {
	ctx, db := connect()
	defer db.Close()

	qry := fmt.Sprintf(`DELETE FROM invoices WHERE fname = '%s'`, fname)
	_, err := db.Exec(ctx, qry)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query or row processing error: %v\n", err)
		os.Exit(1)
	}

	return ReadOp()
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
