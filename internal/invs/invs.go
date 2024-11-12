package invs

import (
	"strings"

	"github.com/ScriptMang/conch/internal/accts"
	"github.com/ScriptMang/conch/internal/bikeshop"
	"github.com/ScriptMang/conch/internal/fields"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type Invoice struct {
	ID       int     `json:"id,omitempty" form:"id,omitempty"`
	UserID   int     `json:"user_id" form:"user_id"`
	Product  string  `json:"product" form:"product"`
	Category string  `json:"category" form:"category"`
	Price    float32 `json:"price" form:"price"`
	Quantity int     `json:"quantity" form:"quantity"`
}

type Invoices []*Invoice

// takes an invoice and throws an error for any field with an invalid input
func (inv *Invoice) validateAllFields(user accts.Users) fields.GrammarError {
	// check for empty fields: for all the fields
	textFields := map[string]*string{
		"Fname":    &user.Fname,
		"Lname":    &user.Lname,
		"Category": &inv.Category,
		"Product":  &inv.Product,
		"Address":  &user.Address,
	}
	var fieldErr fields.GrammarError
	for field, val := range textFields {
		fields.CheckGrammar(field, val, &fieldErr)
	}

	// check for negative values:  price and quantity
	if inv.Price == 0.00 {
		fieldErr.AddMsg(fields.BadRequest, "Error: Price can't be zero")
	} else if inv.Price < 0.00 {
		fieldErr.AddMsg(fields.BadRequest, "Error: The price can't be negative")
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}

	if inv.Quantity == 0 {
		fieldErr.AddMsg(fields.BadRequest, "Error: Quantity can't be zero")
	} else if inv.Quantity < 0 {
		fieldErr.AddMsg(fields.BadRequest, "Error: The quantity can't be negative")
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}
	return fieldErr
}

func InsertOp(usr accts.Users, inv Invoice) ([]*Invoice, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
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

		switch {
		case strings.Contains(qryError, "numeric field overflow"):
			// fmt.Printf("ReadInvoicesByUserID funct: error username doesn't exist\n")
			fieldErr.AddMsg(fields.BadRequest,
				"numeric field overflow, provide a value between 1.00 - 999.99")
		case strings.Contains(qryError, "greater than maximum value for int4"):
			// fmt.Printf("ReadInvoicesByUserID funct: error invoice with specified id doesn't exist\n")
			fieldErr.AddMsg(fields.BadRequest,
				"integer overflow, value must be between 1 - 2147483647")
		case strings.Contains(qryError, "value too long for type character varying"):
			fieldErr.AddMsg(fields.BadRequest, "varchar too long, use varchar length between 1-255")
		default:
			fieldErr.AddMsg(fields.BadRequest, qryError)
		}

	}
	invs = append(invs, &insertedInv)

	return invs, fieldErr
}

// // returns all the invoices in the database a slice []*Invoice
func ReadInvoices() ([]*Invoice, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var invs Invoices
	fieldErr := fields.GrammarError{}
	rows, _ := db.Query(ctx, `SELECT * FROM invoices`)
	err := pgxscan.ScanAll(&invs, rows)
	// fmt.Printf("So Far no errs in ReadInvoices\n")
	if err != nil {
		errMsg := err.Error()
		// fmt.Printf("Houston there's an err in ReadInvoices\n")
		if strings.Contains(errMsg, "failed to connect to `user=username") {
			fieldErr.ErrMsgs = nil
			fieldErr.AddMsg(fields.BadRequest,
				"Error: failed to connect to database, username doesn't exist")
		}
		// fmt.Printf("ReadOp List: %s\n", fieldErr.ErrMsgs)
	}

	return invs, fieldErr
}

func ReadInvoicesByUserID(id int) ([]*Invoice, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var invoices []*Invoice
	_, fieldErr := ReadInvoices()

	if fieldErr.ErrMsgs != nil &&
		strings.Contains(fieldErr.ErrMsgs[0], "failed to connect to `user=username") {
		// log.Printf("ReadInvoicesByUserID funct: Error: username doesn't exist")
		return nil, fieldErr
	}

	rows, _ := db.Query(ctx, `SELECT * FROM invoices WHERE user_id = $1`, id)
	err := pgxscan.ScanAll(&invoices, rows)

	if len(invoices) == 0 {
		// log.Println("Err: No Rows were Found for the Specified User")
		fieldErr.AddMsg(fields.ResourceNotFound, "Resource Not Found: user with specified id doesn't exist")
		return nil, fieldErr
	}

	if err != nil {
		// log.Println("Found an Error Iterating in Getting All the Invoices for the Specified User")
		fieldErr.AddMsg(fields.BadRequest, err.Error())
		return nil, fieldErr
	}

	return invoices, fieldErr
}

// // return the invoice given the user and invoice id
// // if the ids don't exist it returns an error
func ReadInvoiceByUserID(userID, invID int) ([]*Invoice, fields.GrammarError) {
	ctx, db := bikeshop.Connect()
	defer db.Close()

	var invoices []*Invoice
	users, fieldErr := accts.ReadUserByID(userID)

	if fieldErr.ErrMsgs != nil {
		// log.Printf("ReadInvoicesByUserID funct: error username doesn't exist")
		return invoices, fieldErr
	}

	if len(users) == 0 {
		fieldErr.AddMsg(fields.ResourceNotFound, "Resource Not Found: user with specified id doesn't exist")
		return invoices, fieldErr
	}

	rows, _ := db.Query(ctx, `SELECT * FROM invoices WHERE user_id = $1 and id = $2`, userID, invID)

	err := pgxscan.ScanAll(&invoices, rows)

	if len(invoices) == 0 {
		// log.Println("Err: No Rows were Found for the Specified User")
		fieldErr.AddMsg(fields.ResourceNotFound, "Resource Not Found: invoice with specified id doesn't exist")
		return nil, fieldErr
	}

	if err != nil {
		// log.Println("Found an Error Iterating in Getting All the Invoices for the Specified User")
		fieldErr.AddMsg(fields.BadRequest, err.Error())
		return nil, fieldErr
	}

	return invoices, fieldErr
}

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
