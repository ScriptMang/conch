package main

import (
	"fmt"
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
		data := ""
		invs := db.ReadOp()
		for _, inv := range invs {
			str1 := fmt.Sprintf(`"fname": "%s", "lname": "%s", "product": "%s", `, inv.Fname, inv.Lname, inv.Product)
			str2 := fmt.Sprintf(`"price": %.2f, "quantity": %d, "category": "%s", `, inv.Price, inv.Quantity, inv.Category)
			str3 := fmt.Sprintf(`"shipping": "%s"`, inv.Shipping)
			data += fmt.Sprintf(`{` + str1 + str2 + str3 + `},`)
		}
		c.String(http.StatusOK, data)
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

		data := ""
		invs := db.UpdateOp(inv)
		for _, inv2 := range invs {
			str1 := fmt.Sprintf(`"fname": "%s", "lname": "%s", "product": "%s", `, inv2.Fname, inv2.Lname, inv2.Product)
			str2 := fmt.Sprintf(`"price": %.2f, "quantity": %d, "category": "%s", `, inv2.Price, inv2.Quantity, inv2.Category)
			str3 := fmt.Sprintf(`"shipping": "%s"`, inv2.Shipping)
			data += fmt.Sprintf(`{` + str1 + str2 + str3 + `},`)
		}
		c.String(http.StatusOK, data)
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

func main() {
	r := setRouter()
	r = readData(r)

	r = addInvoice(r)
	r = show(r)

	r = updateEntry(r)
	r = showUpdate(r)

	r.Run()
}
