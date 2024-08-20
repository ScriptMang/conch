package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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
	r.POST("/crud1", func(c *gin.Context) {

		var invs db.Invoice
		err := c.ShouldBind(&invs)
		if err != nil {
			log.Fatalf("Error Binding: %v\n", err)
		}

		val := db.InsertOp(invs)
		c.JSON(http.StatusCreated, val)
	})
	return r
}

// reads the tablerows from the database
func readData(r *gin.Engine) *gin.Engine {
	r.GET("/crud2/invoices", func(c *gin.Context) {

		invs := db.ReadInvoices()
		c.JSON(http.StatusOK, invs)
	})
	return r
}

// read a tablerow based on id
func readDataById(r *gin.Engine) *gin.Engine {
	r.GET("/crud2/invoice/:id", func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %v can't converted to an integer\n", err)
		}

		inv := db.ReadInvoiceByID(id)
		c.JSON(http.StatusOK, inv)
	})
	return r
}

// updates an invoice entry by id
func updateEntry(r *gin.Engine) *gin.Engine {
	r.PUT("/crud3/invoice/:id", func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error %v can't converted to an integer\n", err)
		}

		var inv db.Invoice
		if err = c.ShouldBind(&inv); err != nil {
			log.Fatalf("Error Binding: %v\n", err)
		}

		invs := db.UpdateInvoice(inv, id)
		c.JSON(201, invs)
	})
	return r
}

// shows that the invoice entry has been updated
func showDelete(r *gin.Engine) *gin.Engine {
	r.POST("/crud4", func(c *gin.Context) {

		var inv db.Invoice
		if err := c.ShouldBind(&inv); err != nil {
			log.Fatalf("Error Binding: %v\n", err)
		}

		invs := db.DeleteOp(inv)
		c.String(http.StatusOK, invs.Json())
	})
	return r
}

// renders the form page that's needed to Delete an invoice
func deleteEntry(r *gin.Engine) *gin.Engine {
	r.GET("/crud4", func(c *gin.Context) {

		invs := db.Invoices(db.ReadInvoices())

		//generates html option-tags with invoice values
		tmpl := ""
		for _, inv := range invs {
			tmpl += fmt.Sprintf(`<option value='%s'>%s</option>`, inv.Fname, *inv)
		}

		c.HTML(http.StatusOK, "crud4.tmpl", gin.H{
			"title":   "Crud4",
			"details": "Delete an Existing Entry",
			"options": template.HTML(tmpl),
		})

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
	r = showDelete(r)

	r.Run()
}
