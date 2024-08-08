package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	db "github.com/ScriptMang/conch/internal/bikeshop"
	"github.com/gin-gonic/gin"
)

// configs gin router and renders index-page
func setRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Conch Index-Page",
			"crud1": "Create an invoice",
			"crud2": "Print the invoices table",
			"crud3": "Update an existing invoice",
			"crud4": "Delete an existing invoice",
		})
	})
	return r
}

// Takes Post request data of the types: url-encoded or json
// and binds it, to the struct 'invs'.
// When passed to insert-op its used as a bridge
// to add a new invoice.
func show(r *gin.Engine) *gin.Engine {
	r.POST("/crud1", func(c *gin.Context) {

		var invs db.Invoice
		if err := c.ShouldBind(&invs); err != nil {
			log.Fatalf("Error Binding: %v\n", err)
		}

		db.InsertOp(invs)
		c.JSON(http.StatusOK, invs)
	})
	return r
}

// renders the form page that's needed to create an invoice
func addInvoice(r *gin.Engine) *gin.Engine {
	r.GET("/crud1", func(c *gin.Context) {
		c.HTML(http.StatusOK, "crud1.tmpl", gin.H{
			"title":   "Crud1",
			"details": "Add Invoice",
		})
	})
	return r
}

// reads the tablerows from the database
func readData(r *gin.Engine) *gin.Engine {
	r.GET("/crud2", func(c *gin.Context) {
		invs := db.Invoices(db.ReadOp())
		c.String(http.StatusOK, invs.Json())
	})
	return r
}

// shows that the invoice entry has been updated
func showUpdate(r *gin.Engine) *gin.Engine {
	r.POST("/crud3", func(c *gin.Context) {

		var inv db.Invoice
		if err := c.ShouldBind(&inv); err != nil {
			log.Fatalf("Error Binding: %v\n", err)
		}

		invs := db.Invoices(db.UpdateOp(inv))
		c.String(http.StatusOK, invs.Json())
	})
	return r
}

// renders the form page that's needed to update an invoice
func updateEntry(r *gin.Engine) *gin.Engine {
	r.GET("/crud3", func(c *gin.Context) {
		c.HTML(http.StatusOK, "crud3.tmpl", gin.H{
			"title":   "Crud3",
			"details": "Update an Existing Entry",
		})

	})
	return r
}

// shows that the invoice entry has been updated
// func showDelete(r *gin.Engine) *gin.Engine {
// 	r.POST("/crud4", func(c *gin.Context) {

// 		var inv db.Invoice
// 		if err := c.ShouldBind(&inv); err != nil {
// 			log.Fatalf("Error Binding: %v\n", err)
// 		}

// 		invs := db.Invoices(db.DeleteOp(inv))
// 		c.String(http.StatusOK, invs.Json())
// 	})
// 	return r

// renders the form page that's needed to Delete an invoice
func deleteEntry(r *gin.Engine) *gin.Engine {
	r.GET("/crud4", func(c *gin.Context) {

		invs := db.Invoices(db.ReadOp())

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

	r = addInvoice(r)
	r = show(r)

	r = updateEntry(r)
	r = showUpdate(r)

	r = deleteEntry(r)

	r.Run()
}
