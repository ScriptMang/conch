package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	db "github.com/ScriptMang/conch/internal/bikeshop"
	"github.com/gin-gonic/gin"
)

type respBodyData struct {
	Invs     []*db.Invoice
	FieldErr db.GrammarError
}

type resultingInv struct {
	ID       int
	Fname    string
	Lname    string
	Product  string
	Price    json.Number
	Quantity int
	Category string
	Shipping string
}

var code int //httpstatuscode
const statusOK = 200
const statusCreated = 201

// configs gin router and renders index-page
func setRouter() *gin.Engine {
	r := gin.Default()
	return r
}

// assigns an int to an error message that's meant to be modified
func chosenErrorMsg(errMsg string) int {
	// assign the val of 1 for json wrong datatype error
	if strings.Contains(errMsg, "json: cannot unmarshal") {
		return 1
	}
	// assign the val of 2 for incomplete value for key-value pair
	if strings.Contains(errMsg, "looking for beginning of value") {
		return 2
	}
	return 0
}

func editErrMsg(orig, match, edit string) string {
	return strings.Replace(orig, match, edit, 1)
}

// binds an empty invoice to client's data in the response body
// returns the given invoice and an invoice error
func validateInvoiceBinding(c *gin.Context, rqstData *respBodyData) (db.Invoice, bool) {
	var inv db.Invoice
	bindingErr := c.ShouldBind(&inv)

	if bindingErr == nil {
		return inv, true
	}

	err := bindingErr.Error()
	var editedErrMsg string
	errMsgChoice := chosenErrorMsg(err)
	switch errMsgChoice {
	case 1:
		edit1 := editErrMsg(err, "json: cannot unmarshal", "Binding Error:")
		edit2 := editErrMsg(edit1, "into Go struct field Invoice.", "")
		edit3 := editErrMsg(edit2, "of type", "takes a")
		edit4 := edit3 + " not a "
		var temp string
		wordLst := strings.Split(edit4, " ")
		if wordLst[2] == "string" || wordLst[2] == "number" {
			temp = wordLst[2]
			wordLst[2] = ""
		}
		edit5 := strings.Join(wordLst, " ") + temp
		editedErrMsg = editErrMsg(edit5, "  ", " ")
		rqstData.FieldErr.AddMsg(db.BadRequest, editedErrMsg)
	case 2:
		edit := editErrMsg(err, "invalid", "Error: invalid")
		editedErrMsg = editErrMsg(edit,
			"' looking for beginning of value",
			"', value must be wrapped in double quotes")
		rqstData.FieldErr.AddMsg(db.BadRequest, editedErrMsg)
	default:
		editedErrMsg = editErrMsg(err, "invalid", "Error: invalid")
		rqstData.FieldErr.AddMsg(db.BadRequest, editedErrMsg)
	}

	c.AbortWithStatusJSON(db.ErrorCode, rqstData.FieldErr)
	return inv, false
}

func validateRouteID(c *gin.Context, rqstData *respBodyData) int {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rqstData.FieldErr.AddMsg(db.BadRequest, "Bad Request: id can't be converted to an integer")
		sendResponse(c, rqstData)
	}
	return id
}

// serialize Invoice or GrammarError as json to response body
func sendResponse(c *gin.Context, rqstData *respBodyData) {
	invs := rqstData.Invs
	fieldErr := rqstData.FieldErr
	switch {
	case fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "":
		c.JSON(db.ErrorCode, fieldErr)
	default:

		var rsltInv resultingInv
		var invLst []resultingInv

		// change price of type float to string
		// add it to resultingInv struct then invLst
		for _, val := range invs {
			rsltInv.ID = val.ID
			rsltInv.Fname = val.Fname
			rsltInv.Lname = val.Lname
			rsltInv.Product = val.Product
			rsltInv.Price = json.Number(fmt.Sprintf("%.2f", val.Price))
			rsltInv.Quantity = val.Quantity
			rsltInv.Category = val.Category
			rsltInv.Shipping = val.Shipping
			invLst = append(invLst, rsltInv)
		}

		c.JSON(code, invLst)
	}
}

// post request to create user account
func createAcct(r *gin.Engine) *gin.Engine {
	r.POST("/create/Account", func(c *gin.Context) {
		var acct db.Account
		var acctErr db.AuthError
		var fieldErr db.GrammarError
		err := c.ShouldBind(&acct)
		if err != nil {
			acctErr.AddMsg(db.BadRequest,
				"Binding Error: failed to bind fields to account object, mismatched data-types")
			c.JSON(db.ErrorCode, acctErr)
			return
		}

		// validate account info
		db.ValidateAccount(&acct, &fieldErr)
		if len(fieldErr.ErrMsgs) > 0 {
			c.JSON(db.ErrorCode, fieldErr)
			return
		}
		// encrypt password
		acct.Password, err = db.EncryptPassword(acct.Password)
		if err != nil {
			acctErr.AddMsg(db.BadRequest,
				"Hashing Error: password longer than 72 bytes, can't hash")
		}

		// send response back
		errMsgSize := len(acctErr.ErrMsgs)
		switch {
		case errMsgSize > 0:
			c.JSON(db.ErrorCode, acctErr)
		default:
			c.JSON(statusOK, acct)
		}

		//log.Println("Account: ", acct)
	})
	return r
}

// binds json data to an invoice and insert its to the database
func addInvoice(r *gin.Engine) *gin.Engine {
	r.POST("/invoices/", func(c *gin.Context) {
		var inv db.Invoice
		var rqstData respBodyData
		var bindingOk bool
		inv, bindingOk = validateInvoiceBinding(c, &rqstData)
		if bindingOk {
			rqstData.Invs, rqstData.FieldErr = db.InsertOp(inv)
			code = statusCreated
			sendResponse(c, &rqstData)
		}
	})
	return r
}

// reads the tablerows from the database
func readData(r *gin.Engine) *gin.Engine {
	r.GET("/invoices", func(c *gin.Context) {
		var rqstData respBodyData
		rqstData.Invs, rqstData.FieldErr = db.ReadInvoices()
		code = statusOK
		sendResponse(c, &rqstData)
	})
	return r
}

// read a tablerow based on id
func readDataById(r *gin.Engine) *gin.Engine {
	r.GET("/invoice/:id", func(c *gin.Context) {
		var rqstData respBodyData
		id := validateRouteID(c, &rqstData)
		if id != 0 {
			rqstData.Invs, rqstData.FieldErr = db.ReadInvoiceByID(id)
			code = statusOK
			sendResponse(c, &rqstData)
		}
	})
	return r
}

// updates an invoice entry by id
// require the user to pass the entire invoice
// to change any field
func updateEntry(r *gin.Engine) *gin.Engine {
	r.PUT("/invoice/:id", func(c *gin.Context) {
		var inv db.Invoice
		var bindingOk bool
		var rqstData respBodyData
		id := validateRouteID(c, &rqstData)
		if id != 0 {
			inv, bindingOk = validateInvoiceBinding(c, &rqstData)
			if bindingOk {
				rqstData.Invs, rqstData.FieldErr = db.UpdateInvoice(inv, id)
				code = statusOK
				sendResponse(c, &rqstData)
			}
		}
	})
	return r
}

// similar to the updateEntry except you don't have
// to pass all the fields in a invoice to update a field
func patchEntry(r *gin.Engine) *gin.Engine {
	r.PATCH("/invoice/:id", func(c *gin.Context) {
		var inv db.Invoice
		var bindingOk bool
		var rqstData respBodyData
		id := validateRouteID(c, &rqstData)
		if id != 0 {
			inv, bindingOk = validateInvoiceBinding(c, &rqstData)
			if bindingOk {
				rqstData.Invs, rqstData.FieldErr = db.PatchInvoice(inv, id)
				code = statusOK
				sendResponse(c, &rqstData)
			}
		}
	})
	return r
}

// deletes an invoice entry based on id
func deleteEntry(r *gin.Engine) *gin.Engine {
	r.DELETE("/invoice/:id", func(c *gin.Context) {
		var rqstData respBodyData
		id := validateRouteID(c, &rqstData)
		if id != 0 {
			rqstData.Invs, rqstData.FieldErr = db.DeleteInvoice(id)
			code = statusOK
			sendResponse(c, &rqstData)
		}
	})
	return r
}

func main() {
	r := setRouter()

	r = createAcct(r)
	r = readData(r)
	r = readDataById(r)
	r = addInvoice(r)
	r = updateEntry(r)
	r = patchEntry(r)
	r = deleteEntry(r)

	r.Run()
}
