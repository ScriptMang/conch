package main

import (
	"net/http"
	"strconv"

	db "github.com/ScriptMang/conch/internal/bikeshop"
	"github.com/gin-gonic/gin"
)

// configs gin router and renders index-page
func setRouter() *gin.Engine {
	r := gin.Default()
	return r
}

// Takes Post request data of the types: url-encoded or json
// and binds it, to the struct 'invs'.
// When passed to insert-op its used as a bridge
// to add a new invoice.
func addInvoice(r *gin.Engine) *gin.Engine {
	r.POST("/crud1/invoices/", func(c *gin.Context) {
		var invs db.Invoice
		var fieldErr db.InvoiceError
		err := c.BindJSON(&invs)
		if err != nil {
			fieldErr.HttpStatusCode = 415
			fieldErr.Msg = append(fieldErr.Msg, "Failed to bind invoice, request only takes JSON data")
		}

		inv, fieldErr := db.InsertOp(invs)
		switch {
		case len(fieldErr.Msg) > 0:
			c.JSON(fieldErr.HttpStatusCode, fieldErr)
		default:
			c.JSON(http.StatusCreated, inv)
		}
	})
	return r
}

// reads the tablerows from the database
func readData(r *gin.Engine) *gin.Engine {
	r.GET("/crud2/invoices", func(c *gin.Context) {
		invs, tableErr := db.ReadInvoices()

		switch {
		case len(tableErr.Msg) > 0:
			c.JSON(tableErr.HttpStatusCode, tableErr)
		default:
			c.JSON(http.StatusOK, invs)
		}
	})
	return r
}

// read a tablerow based on id
func readDataById(r *gin.Engine) *gin.Engine {
	r.GET("/crud2/invoice/:id", func(c *gin.Context) {
		var inv db.Invoice
		var fieldErr db.InvoiceError
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fieldErr.ContentType = "application/json"
			fieldErr.HttpStatusCode = 400
			fieldErr.Msg = append(fieldErr.Msg, "Bad Request: id can't converted to an integer")
		} else {
			inv, fieldErr = db.ReadInvoiceByID(id)
		}

		switch {
		case len(fieldErr.Msg) > 0:
			c.JSON(fieldErr.HttpStatusCode, fieldErr)
		default:
			c.JSON(http.StatusOK, inv)
		}
	})
	return r
}

// updates an invoice entry by id
func updateEntry(r *gin.Engine) *gin.Engine {
	r.PUT("/crud3/invoice/:id", func(c *gin.Context) {

		var inv, inv2 db.Invoice
		var fieldErr db.InvoiceError
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fieldErr.ContentType = "application/json"
			fieldErr.HttpStatusCode = 400
			fieldErr.Msg = append(fieldErr.Msg, "Bad Request: id can't converted to an integer")
		} else {
			bindingErr := c.BindJSON(&inv)
			if bindingErr != nil {
				fieldErr.HttpStatusCode = 415
				fieldErr.Msg = append(fieldErr.Msg, "Failed to bind invoice, request only takes JSON data")
			}
			inv2, fieldErr = db.UpdateInvoice(inv, id)
		}

		switch {
		case len(fieldErr.Msg) > 0:
			c.JSON(fieldErr.HttpStatusCode, fieldErr)
		default:
			c.JSON(http.StatusCreated, inv2)
		}
	})
	return r
}

// deletes an invoice entry based on id
func deleteEntry(r *gin.Engine) *gin.Engine {
	r.DELETE("/crud4/invoice/:id", func(c *gin.Context) {

		var inv db.Invoice
		var fieldErr db.InvoiceError
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fieldErr.ContentType = "application/json"
			fieldErr.HttpStatusCode = 400
			fieldErr.Msg = append(fieldErr.Msg, "Bad Request: id can't converted to an integer")
		} else {
			inv, fieldErr = db.DeleteInvoice(id)
		}

		switch {
		case len(fieldErr.Msg) > 0:
			c.JSON(fieldErr.HttpStatusCode, fieldErr)
		default:
			c.JSON(http.StatusOK, inv)
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
