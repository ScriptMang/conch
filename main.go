package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/ScriptMang/conch/internal/accts"
	"github.com/ScriptMang/conch/internal/fields"
	"github.com/ScriptMang/conch/internal/invs"
	"github.com/gin-gonic/gin"
)

type respBodyData struct {
	Invs        []*invs.Invoice
	UsrContacts []*accts.UserContacts
	FieldErr    fields.GrammarError
}

type Order struct {
	ID       int         `json:"id"`
	UserID   int         `json:"user_id"`
	Fname    string      `json:"fname"`
	Lname    string      `json:"lname"`
	Product  string      `json:"product"`
	Price    json.Number `json:"price"`
	Quantity int         `json:"quantity"`
	Category string      `json:"category"`
	Address  string      `json:"address"`
}

var code int //httpstatuscode
const statusOK = 200
const statusCreated = 201

var btokens []accts.Tokens // bearer token

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
// func validateRouteUserID(c *gin.Context, rqstData *respBodyData) int {
// 	id, err := strconv.Atoi(c.Param("usr_id"))
// 	if err != nil {
// 		rqstData.FieldErr.AddMsg(fields.BadRequest, "Bad Request: user id can't be converted to an integer")
// 	}
// 	return id
// }

// validates the invoice id route parameter
func validateRouteInvID(c *gin.Context, rqstData *respBodyData) int {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rqstData.FieldErr.AddMsg(fields.BadRequest, "Bad Request: invoice id can't be converted to an integer")
	}
	return id
}

// serialize Invoice or GrammarError as json to response body
func sendResponse(c *gin.Context, rqstData *respBodyData) {
	invs := rqstData.Invs
	usrContacts := rqstData.UsrContacts
	fieldErr := rqstData.FieldErr
	switch {
	case fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "":
		c.JSON(fields.ErrorCode, fieldErr)
	default:

		var receipt Order
		var receipts []Order

		// change price of type float to string
		// add it to resultingInv struct then invLst
		for _, usrContact := range usrContacts {
			for _, inv := range invs {
				if inv.UserID != usrContact.ID {
					continue
				}
				receipt.ID = inv.ID
				receipt.UserID = inv.UserID
				receipt.Product = inv.Product
				receipt.Price = json.Number(fmt.Sprintf("%.2f", inv.Price))
				receipt.Quantity = inv.Quantity
				receipt.Category = inv.Category
				receipt.Fname = usrContact.Fname
				receipt.Lname = usrContact.Lname
				receipt.Address = usrContact.Address
				receipts = append(receipts, receipt)
			}
		}
		c.JSON(code, receipts)
	}
}

// post request to create user account
func createAcct(r *gin.Engine) *gin.Engine {
	r.POST("/users", func(c *gin.Context) {
		var acct accts.Account
		var acctErr fields.GrammarError
		var acctStatus *accts.Registered
		err := c.ShouldBind(&acct)
		if err != nil {
			acctErr.AddMsg(fields.BadRequest,
				"Binding Error: failed to bind fields to account object, mismatched data-types")
			c.JSON(fields.ErrorCode, acctErr)
			return
		}

		// validate account info
		acctStatus, acctErr = accts.AddAccount(&acct)
		// send response back
		errMsgSize := len(acctErr.ErrMsgs)
		switch {
		case errMsgSize > 0:
			c.JSON(fields.ErrorCode, acctErr)
		default:
			c.JSON(statusOK, *acctStatus)
		}
	})
	return r
}

func logIn(c *gin.Context) {
	var username, ok = c.MustGet(gin.AuthUserKey).(string)
	if !ok {
		log.Fatalf("Couldn't Get user's username from AUTHUserKey\n")
	}
	var fieldErr fields.GrammarError
	token := accts.GenerateToken(username, &fieldErr)
	if fieldErr.ErrMsgs != nil {
		c.JSON(accts.BadRequest, gin.H{
			"TokenError": fieldErr.ErrMsgs,
		})
		return
	}

	btokens = append(btokens, token)
	c.JSON(http.StatusAccepted, gin.H{
		"token": string(token.Token),
	})
}

func logOut(c *gin.Context) {

	if c.Keys["isAuthorized"] == false {
		return
	}

	// verify that the userID assigned to the token matches the route's userID
	if c.Keys["rqstTokenUserID"] == 0 {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var fieldErr fields.GrammarError
	userID := c.Keys["rqstTokenUserID"].(int)
	username := accts.ReadUsernameByID(userID, &fieldErr)
	if fieldErr.ErrMsgs != nil {
		c.JSON(accts.BadRequest, gin.H{
			"Error": fieldErr.ErrMsgs,
		})
		return
	}

	accts.LogOut(userID, &fieldErr)
	if fieldErr.ErrMsgs != nil {
		c.JSON(accts.BadRequest, gin.H{
			"TokenError": fieldErr.ErrMsgs,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("User: %s has logged out", username[0].Username),
	})

}

func deleteAcct(c *gin.Context) {

	if c.Keys["isAuthorized"] == false {
		return
	}

	var rqstData respBodyData
	var rmvUser []*accts.Usernames

	// verify that the userID assigned to the token matches the route's userID
	if c.Keys["rqstTokenUserID"] == 0 {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var fieldErr fields.GrammarError
	userID := c.Keys["rqstTokenUserID"].(int)
	user := accts.ReadUsernameByID(userID, &fieldErr)
	if fieldErr.ErrMsgs != nil {
		c.JSON(accts.BadRequest, gin.H{
			"Error": fieldErr.ErrMsgs,
		})
		return
	}

	rmvUser, rqstData.FieldErr = accts.DeleteAcct(*user[0])
	if rqstData.FieldErr.ErrMsgs != nil {
		sendResponse(c, &rqstData)
		return
	}

	code = statusOK
	c.JSON(code, gin.H{
		"message": fmt.Sprintf("User: %s has been deleted", rmvUser[0].Username),
	})
}

func protectData(c *gin.Context) {
	c.Keys = make(map[string]any)
	bToken := c.Request.Header.Get("Authorization")
	if bToken == "" {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	// get the userid for the inputted token
	var fieldErr fields.GrammarError
	rqstToken := strings.Split(bToken, " ")[1]
	rqstTokenUserID := accts.ReadUserIDByToken(rqstToken, &fieldErr)

	// return the error retrieving the token
	if fieldErr.ErrMsgs != nil {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"TokenError": fieldErr.ErrMsgs,
		})
		return
	}
	c.Keys["rqstTokenUserID"] = rqstTokenUserID

	for _, token := range btokens {
		dbTknUserID := accts.ReadUserIDByToken(string(token.Token), &fieldErr)
		if dbTknUserID == rqstTokenUserID {
			c.Keys["isAuthorized"] = true
			return
		}
	}

	c.Keys["isAuthorized"] = false
	c.JSON(http.StatusUnauthorized, gin.H{
		"message": "unauthorized",
	})
}

// binds json data to an invoice and insert its to the database
func addInvoice(c *gin.Context) {
	if c.Keys["isAuthorized"] == false {
		return
	}
	var inv invs.Invoice
	var rqstData respBodyData
	var bindingOk bool
	inv, bindingOk = validateInvoiceBinding(c, &rqstData)
	// fmt.Printf("Invoice in addInvoice funct is: %+v\n", inv)

	// invoice binding must match along with their userID
	if c.Keys["rqstTokenUserID"] == 0 {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	if bindingOk {
		var fieldErr fields.GrammarError
		inv.UserID = c.Keys["rqstTokenUserID"].(int)
		rqstData.Invs, fieldErr = invs.InsertOp(inv)
		// fmt.Printf("Invoice after InsertOP is: %+v\n", *rqstData.Invs[0])
		// fmt.Printf("FieldErrs after InsertOP is: %v\n", fieldErr.ErrMsgs)
		if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
			rqstData.FieldErr = fieldErr
			sendResponse(c, &rqstData)
			return
		}
		code = statusCreated
		inv2 := *rqstData.Invs[0]
		c.JSON(code, gin.H{
			"id":       inv2.ID,
			"product":  inv2.Product,
			"category": inv2.Category,
			"price":    inv2.Price,
			"quantity": inv2.Quantity,
		})
	}
}

// returns the list of all users
func readUserData(c *gin.Context) {
	if c.Keys["isAuthorized"] == false {
		return
	}
	var rqstData respBodyData
	rqstData.UsrContacts, rqstData.FieldErr = accts.ReadUserContact()
	fieldErr := rqstData.FieldErr
	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		// log.Printf("Error in ReadUserData funct: %v\n", fieldErr.ErrMsgs)
		sendResponse(c, &rqstData)
		return
	}
	code = statusOK
	c.JSON(code, rqstData.UsrContacts)
}

// returns a user given its id
func readUserDataByID(c *gin.Context) {
	// verify first whether the token exist in the databse
	if c.Keys["isAuthorized"] == false {
		return
	}

	// verify that the userID assigned to the token matches the route's userID
	if c.Keys["rqstTokenUserID"] == 0 {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var rqstData respBodyData
	id := c.Keys["rqstTokenUserID"].(int)
	rqstData.UsrContacts, rqstData.FieldErr = accts.ReadUserContactByID(id)
	fieldErr := rqstData.FieldErr
	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		// fmt.Printf("readUserInovices funct: error is %s\n", fieldErr.ErrMsgs[0])
		sendResponse(c, &rqstData)
		return
	}
	code = statusOK
	c.JSON(code, *rqstData.UsrContacts[0])
}

// // returns all the invoices within the database
func readInvoiceData(c *gin.Context) {
	if c.Keys["isAuthorized"] == false {
		return
	}
	var rqstData respBodyData
	rqstData.Invs, rqstData.FieldErr = invs.ReadInvoices()
	fieldErr := rqstData.FieldErr
	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		sendResponse(c, &rqstData)
		return
	}
	code = statusOK
	c.JSON(code, rqstData.Invs)
}

// returns all the invoices for a given user
func readUserInvoices(c *gin.Context) {
	if c.Keys["isAuthorized"] == false {
		return
	}

	// verify that the userID assigned to the token matches the route's userID
	if c.Keys["rqstTokenUserID"] == 0 {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	var rqstData respBodyData
	id := c.Keys["rqstTokenUserID"].(int)
	rqstData.Invs, rqstData.FieldErr = invs.ReadInvoicesByUserID(id)
	fieldErr := rqstData.FieldErr
	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		// fmt.Printf("readUserInovices funct: error is %s\n", fieldErr.ErrMsgs[0])
		sendResponse(c, &rqstData)
		return
	}
	code = statusOK
	c.JSON(code, rqstData.Invs)
}

// returns a specific invoice for a specific user
// given the user id and invoice id
func readUserInvoiceByID(c *gin.Context) {
	if c.Keys["isAuthorized"] == false {
		return
	}

	var rqstData respBodyData
	invID := validateRouteInvID(c, &rqstData)

	// verify that the userID assigned to the token matches the route's userID
	if c.Keys["rqstTokenUserID"] == 0 {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	userID := c.Keys["rqstTokenUserID"].(int)
	rqstData.Invs, rqstData.FieldErr = invs.ReadInvoiceByUserID(userID, invID)
	fieldErr := rqstData.FieldErr
	if fieldErr.ErrMsgs != nil && fieldErr.ErrMsgs[0] != "" {
		// fmt.Printf("readUserInovices funct: error is %s\n", fieldErr.ErrMsgs[0])
		sendResponse(c, &rqstData)
		return
	}
	code = statusOK
	c.JSON(code, rqstData.Invs)
}

// updates an invoice entry by id
// require the user to pass the entire invoice
// to change any field
func updateInvoiceEntry(c *gin.Context) {
	if c.Keys["isAuthorized"] == false {
		return
	}
	var inv invs.Invoice
	var bindingOk bool

	var rqstData respBodyData
	invID := validateRouteInvID(c, &rqstData)

	// verify that the userID assigned to the token matches the route's userID
	if c.Keys["rqstTokenUserID"] == 0 {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	inv, bindingOk = validateInvoiceBinding(c, &rqstData)
	if bindingOk {

		userID := c.Keys["rqstTokenUserID"].(int)
		rqstData.Invs, rqstData.FieldErr = invs.UpdateInvoiceByUserID(inv, userID, invID)
		if rqstData.FieldErr.ErrMsgs != nil {
			sendResponse(c, &rqstData)
			return
		}
		code = statusOK
		c.JSON(code, rqstData.Invs)
	}
}

// similar to the updateEntry except you don't have
// to pass all the fields in a invoice to update a field
func patchEntry(c *gin.Context) {
	if c.Keys["isAuthorized"] == false {
		return
	}
	var inv invs.Invoice
	var bindingOk bool

	var rqstData respBodyData
	invID := validateRouteInvID(c, &rqstData)

	// verify that the userID assigned to the token matches the route's userID
	if c.Keys["rqstTokenUserID"] == 0 {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	inv, bindingOk = validateInvoiceBinding(c, &rqstData)
	if bindingOk {
		userID := c.Keys["rqstTokenUserID"].(int)
		rqstData.Invs, rqstData.FieldErr = invs.PatchInvoice(inv, userID, invID)
		if rqstData.FieldErr.ErrMsgs != nil {
			sendResponse(c, &rqstData)
			return
		}
		code = statusOK
		c.JSON(code, rqstData.Invs)
	}
}

// deletes an invoice entry based on id
func deleteInvEntry(c *gin.Context) {
	if c.Keys["isAuthorized"] == false {
		return
	}

	var rqstData respBodyData
	invID := validateRouteInvID(c, &rqstData)

	if c.Keys["rqstTokenUserID"] == 0 {
		c.Keys["isAuthorized"] = false
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	userID := c.Keys["rqstTokenUserID"].(int)
	rqstData.Invs, rqstData.FieldErr = invs.DeleteInvoice(invID, userID)
	if rqstData.FieldErr.ErrMsgs != nil {
		sendResponse(c, &rqstData)
		return
	}
	code = statusOK
	c.JSON(code, rqstData.Invs)
}

func main() {
	r := setRouter()
	r = createAcct(r)

	isHashUnreadable := false
	const hash_unreadable = "Couldn't read password hash"
	pswd1, readErr := accts.ReadHashByID(1)
	if readErr.ErrMsgs != nil {
		isHashUnreadable = true
		fmt.Fprintf(os.Stderr, hash_unreadable+" 1\n")
	}

	pswd2, readErr := accts.ReadHashByID(2)
	if readErr.ErrMsgs != nil {
		isHashUnreadable = true
		fmt.Fprintf(os.Stderr, hash_unreadable+" 2\n")
	}

	if !isHashUnreadable {
		loginRoute := r.Group("/", gin.BasicAuth(gin.Accounts{
			"wrigglyWart56": string(pswd1[0].Password),
			"hypnoTonic05":  string(pswd2[0].Password),
		}))
		{
			loginRoute.POST("/login", logIn)
		}

		userGroup1 := r.Group("/", protectData)
		{
			userGroup1.GET("/users", readUserData)
			userGroup1.GET("/invoices", readInvoiceData)
			userGroup1.DELETE("/users", deleteAcct)
			userGroup1.POST("/logout", logOut)
		}

		createInv := r.Group("/invoices/", protectData)
		{
			createInv.POST("", addInvoice)
		}

		userGroup2 := r.Group("/", protectData)
		{
			userGroup2.GET("/user", readUserDataByID)           // read user by their id
			userGroup2.GET("/user/invoices", readUserInvoices)  // read all the invoices for a user
			userGroup2.GET("/invoice/:id", readUserInvoiceByID) // read a specific invoice from a user
			userGroup2.PUT("/invoice/:id", updateInvoiceEntry)  // updates the entire invoice
			userGroup2.PATCH("/invoice/:id", patchEntry)        // updates any field of an invoice
			userGroup2.DELETE("/invoice/:id", deleteInvEntry)   // deletes a specific invoice
		}
	}

	r.Run()
}
