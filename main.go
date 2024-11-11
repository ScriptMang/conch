package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/ScriptMang/conch/internal/accts"
	"github.com/ScriptMang/conch/internal/fields"
	"github.com/ScriptMang/conch/internal/invs"
	"github.com/gin-gonic/gin"
)

type respBodyData struct {
	Invs     []*invs.Invoice
	Users    []*accts.Users
	FieldErr fields.GrammarError
}

type Order struct {
	ID       int
	UserID   int
	Fname    string
	Lname    string
	Product  string
	Price    json.Number
	Quantity int
	Category string
	Address  string
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
func validateInvoiceBinding(c *gin.Context, rqstData *respBodyData) (invs.Invoice, bool) {
	var inv invs.Invoice
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
		rqstData.FieldErr.AddMsg(fields.BadRequest, editedErrMsg)
	case 2:
		edit := editErrMsg(err, "invalid", "Error: invalid")
		editedErrMsg = editErrMsg(edit,
			"' looking for beginning of value",
			"', value must be wrapped in double quotes")
		rqstData.FieldErr.AddMsg(fields.BadRequest, editedErrMsg)
	default:
		editedErrMsg = editErrMsg(err, "invalid", "Error: invalid")
		rqstData.FieldErr.AddMsg(fields.BadRequest, editedErrMsg)
	}

	c.AbortWithStatusJSON(fields.ErrorCode, rqstData.FieldErr)
	return inv, false
}

// validates the user id route parameter
func validateRouteUserID(c *gin.Context, rqstData *respBodyData) int {
	id, err := strconv.Atoi(c.Param("usr_id"))
	if err != nil {
		rqstData.FieldErr.AddMsg(fields.BadRequest, "Bad Request: id can't be converted to an integer")
		sendResponse(c, rqstData)
	}
	return id
}

// validates the invoice id route parameter
func validateRouteInvID(c *gin.Context, rqstData *respBodyData) int {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rqstData.FieldErr.AddMsg(fields.BadRequest, "Bad Request: id can't be converted to an integer")
		sendResponse(c, rqstData)
	}
	return id
}

// serialize Invoice or GrammarError as json to response body
func sendResponse(c *gin.Context, rqstData *respBodyData) {
	invs := rqstData.Invs
	usrs := rqstData.Users
	fieldErr := rqstData.FieldErr
	switch {
	case fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "":
		c.JSON(fields.ErrorCode, fieldErr)
	default:

		var receipt Order
		var receipts []Order

		// change price of type float to string
		// add it to resultingInv struct then invLst
		for _, usr := range usrs {
			for _, inv := range invs {
				if inv.UserID != usr.ID {
					continue
				}
				receipt.ID = inv.ID
				receipt.UserID = inv.UserID
				receipt.Product = inv.Product
				receipt.Price = json.Number(fmt.Sprintf("%.2f", inv.Price))
				receipt.Quantity = inv.Quantity
				receipt.Category = inv.Category
				receipt.Fname = usr.Fname
				receipt.Lname = usr.Lname
				receipt.Address = usr.Address
				receipts = append(receipts, receipt)
			}
		}
		c.JSON(code, receipts)
	}
}

// post request to create user account
func createAcct(r *gin.Engine) *gin.Engine {
	r.POST("/create/Account", func(c *gin.Context) {
		var acct *accts.Account
		var acctErr fields.GrammarError
		var respData []*accts.Account
		err := c.ShouldBind(&acct)
		if err != nil {
			acctErr.AddMsg(fields.BadRequest,
				"Binding Error: failed to bind fields to account object, mismatched data-types")
			c.JSON(fields.ErrorCode, acctErr)
			return
		}

		// validate account info
		respData, acctErr = accts.AddAccount(acct)

		if len(respData) == 0 {
			fmt.Println("Thats strange, no accounts were added")
		}
		// send response back
		errMsgSize := len(acctErr.ErrMsgs)
		switch {
		case errMsgSize > 0:
			c.JSON(fields.ErrorCode, acctErr)
		default:
			c.JSON(statusOK, *respData[0])
		}

		//log.Println("Account: ", acct)
	})
	return r
}

// binds json data to an invoice and insert its to the database
func addInvoice(r *gin.Engine) *gin.Engine {
	r.POST("/user/:usr_id/invoices/", func(c *gin.Context) {
		var inv invs.Invoice
		var rqstData respBodyData
		var bindingOk bool
		id := validateRouteUserID(c, &rqstData)
		// fmt.Printf("ID in addInvoice funct is: %d\n", id)
		inv, bindingOk = validateInvoiceBinding(c, &rqstData)
		// fmt.Printf("Invoice in addInvoice funct is: %+v\n", inv)
		if id != 0 && bindingOk {
			var fieldErr fields.GrammarError
			rqstData.Users, fieldErr = accts.ReadUserByID(id)
			usrs := rqstData.Users
			// fmt.Printf("User in addInvoice funct is: %+v\n", *usrs[0])
			// fmt.Printf("The length of fieldErrs after ReadUserByID is: %d\n", len(fieldErr.ErrMsgs))
			if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
				rqstData.FieldErr = fieldErr
				sendResponse(c, &rqstData)
				return
			}
			rqstData.Invs, fieldErr = invs.InsertOp(*usrs[0], inv)
			rqstData.FieldErr = fieldErr
			// fmt.Printf("Invoice after InsertOP is: %+v\n", *rqstData.Invs[0])
			// fmt.Printf("FieldErrs after InsertOP is: %v\n", fieldErr.ErrMsgs)
			code = statusCreated
			sendResponse(c, &rqstData)
		}
	})
	return r
}

// returns the list of all users
func readUserData(r *gin.Engine) *gin.Engine {
	r.GET("/users", func(c *gin.Context) {
		var rqstData respBodyData
		rqstData.Users, rqstData.FieldErr = accts.ReadUsers()
		fieldErr := rqstData.FieldErr
		if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
			// log.Printf("Error in ReadUserData funct: %v\n", fieldErr.ErrMsgs)
			sendResponse(c, &rqstData)
			return
		}
		code = statusOK
		c.JSON(code, rqstData.Users)
	})
	return r
}

// returns a user given its id
func readUserDataByID(r *gin.Engine) *gin.Engine {
	r.GET("/user/:usr_id", func(c *gin.Context) {
		var rqstData respBodyData
		id := validateRouteUserID(c, &rqstData)
		var invalidID = rqstData.FieldErr.ErrMsgs
		if invalidID != nil {
			return
		}
		rqstData.Users, rqstData.FieldErr = accts.ReadUserByID(id)
		fieldErr := rqstData.FieldErr
		if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
			// fmt.Printf("readUserInovices funct: error is %s\n", fieldErr.ErrMsgs[0])
			sendResponse(c, &rqstData)
			return
		}
		code = statusOK
		c.JSON(code, *rqstData.Users[0])
	})
	return r
}

// // returns all the invoices within the database
func readInvoiceData(r *gin.Engine) *gin.Engine {
	r.GET("/users/invoices", func(c *gin.Context) {
		var rqstData respBodyData
		rqstData.Invs, rqstData.FieldErr = invs.ReadInvoices()
		fieldErr := rqstData.FieldErr
		if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
			sendResponse(c, &rqstData)
			return
		}
		code = statusOK
		c.JSON(code, rqstData.Invs)
	})
	return r
}

// returns all the invoices for a given user
func readUserInvoices(r *gin.Engine) *gin.Engine {
	r.GET("/user/:usr_id/invoices", func(c *gin.Context) {
		var rqstData respBodyData
		id := validateRouteUserID(c, &rqstData)
		var invalidID = rqstData.FieldErr.ErrMsgs
		if invalidID != nil {
			return
		}
		rqstData.Invs, rqstData.FieldErr = invs.ReadInvoicesByUserID(id)
		fieldErr := rqstData.FieldErr
		if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
			// fmt.Printf("readUserInovices funct: error is %s\n", fieldErr.ErrMsgs[0])
			sendResponse(c, &rqstData)
			return
		}
		code = statusOK
		c.JSON(code, rqstData.Invs)
	})
	return r
}

// returns a specific invoice for a specific user
// given the user id and invoice id
func readUserInvoiceByID(r *gin.Engine) *gin.Engine {
	r.GET("/user/:usr_id/invoice/:id", func(c *gin.Context) {
		var rqstData respBodyData
		userID := validateRouteUserID(c, &rqstData)
		invID := validateRouteInvID(c, &rqstData)
		rqstData.Invs, rqstData.FieldErr = invs.ReadInvoiceByUserID(userID, invID)
		fieldErr := rqstData.FieldErr
		if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
			// fmt.Printf("readUserInovices funct: error is %s\n", fieldErr.ErrMsgs[0])
			sendResponse(c, &rqstData)
			return
		}
		code = statusOK
		c.JSON(code, rqstData.Invs)
	})
	return r
}

// // updates an invoice entry by id
// // require the user to pass the entire invoice
// // to change any field
// func updateEntry(r *gin.Engine) *gin.Engine {
// 	r.PUT("/invoice/:id", func(c *gin.Context) {
// 		var inv db.Invoice
// 		var bindingOk bool
// 		var rqstData respBodyData
// 		id := validateRouteID(c, &rqstData)
// 		if id != 0 {
// 			inv, bindingOk = validateInvoiceBinding(c, &rqstData)
// 			if bindingOk {
// 				rqstData.Invs, rqstData.FieldErr = db.UpdateInvoice(inv, id)
// 				code = statusOK
// 				sendResponse(c, &rqstData)
// 			}
// 		}
// 	})
// 	return r
// }

// // similar to the updateEntry except you don't have
// // to pass all the fields in a invoice to update a field
// func patchEntry(r *gin.Engine) *gin.Engine {
// 	r.PATCH("/invoice/:id", func(c *gin.Context) {
// 		var inv db.Invoice
// 		var bindingOk bool
// 		var rqstData respBodyData
// 		id := validateRouteID(c, &rqstData)
// 		if id != 0 {
// 			inv, bindingOk = validateInvoiceBinding(c, &rqstData)
// 			if bindingOk {
// 				rqstData.Invs, rqstData.FieldErr = db.PatchInvoice(inv, id)
// 				code = statusOK
// 				sendResponse(c, &rqstData)
// 			}
// 		}
// 	})
// 	return r
// }

// // deletes an invoice entry based on id
// func deleteEntry(r *gin.Engine) *gin.Engine {
// 	r.DELETE("/invoice/:id", func(c *gin.Context) {
// 		var rqstData respBodyData
// 		id := validateRouteID(c, &rqstData)
// 		if id != 0 {
// 			rqstData.Invs, rqstData.FieldErr = db.DeleteInvoice(id)
// 			code = statusOK
// 			sendResponse(c, &rqstData)
// 		}
// 	})
// 	return r
// }

func main() {
	r := setRouter()

	r = createAcct(r)
	r = readUserData(r)
	r = readInvoiceData(r)
	r = readUserDataByID(r)
	r = readUserInvoices(r)
	r = readUserInvoiceByID(r)
	r = addInvoice(r)
	// r = updateEntry(r)
	// r = patchEntry(r)
	// r = deleteEntry(r)

	r.Run()
}
