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
	FieldErr db.InvoiceError
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

// configs gin router and renders index-page
func setRouter() *gin.Engine {
	r := gin.Default()
	return r
}

// binds an empty invoice to client's data in the response body
// returns the given invoice and an invoice error
func validateInvoiceBinding(c *gin.Context, rqstData *respBodyData) (db.Invoice, bool) {
	var inv db.Invoice
	bindingErr := c.ShouldBind(&inv)
	if bindingErr != nil {
		err := bindingErr.Error()
		var editedErrMsg string

		if strings.Contains(err, "json: cannot unmarshal") {
			editedErrMsg = strings.Replace(err, "json:", "Error: json", 1)
			rqstData.FieldErr.AddMsg(400, editedErrMsg)
		} else if strings.Contains(err, "looking for beginning of value") {
			editedErrMsg = strings.Replace(err, "invalid", "Error: invalid", 1)
			editedErrMsg = strings.Replace(
				editedErrMsg, "' looking for beginning of value",
				"', value must be wrapped in double quotes", 1)
			rqstData.FieldErr.AddMsg(400, editedErrMsg)
		} else {
			editedErrMsg = strings.Replace(err, "invalid", "Error: invalid", 1)
			rqstData.FieldErr.AddMsg(400, editedErrMsg)
		}

		c.AbortWithStatusJSON(db.ErrorCode, rqstData.FieldErr)
		return inv, false
	}
	return inv, true
}

func validateRouteID(c *gin.Context, rqstData *respBodyData) int {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		rqstData.FieldErr.AddMsg(400, "Bad Request: id can't be converted to an integer")
		sendResponse(c, rqstData)
	}
	return id
}

// serialize Invoice or InvoiceError as json to response body
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

// // binds json data to an invoice and insert its to the database
func addInvoice(r *gin.Engine) *gin.Engine {
	r.POST("/invoices/", func(c *gin.Context) {
		var inv db.Invoice
		var rqstData respBodyData
		var bindingOk bool
		inv, bindingOk = validateInvoiceBinding(c, &rqstData)
		if bindingOk {
			rqstData.Invs, rqstData.FieldErr = db.InsertOp(inv)
			code = 201
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
		code = 200
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
			code = 200
			sendResponse(c, &rqstData)
		}
	})
	return r
}

// updates an invoice entry by id
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
				code = 201
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
			code = 200
			sendResponse(c, &rqstData)
		}
	})
	return r
}

func main() {
	r := setRouter()

	r = readData(r)
	r = readDataById(r)
	r = addInvoice(r)
	r = updateEntry(r)
	r = deleteEntry(r)

	r.Run()
}
