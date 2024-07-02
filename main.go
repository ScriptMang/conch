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
	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Conch Index-Page",
			"crud1": "Insert Table-Row Data",
			"crud2": "Perform a Read Operation",
		})
	})
	return r
}

// sends form data as json Post request
func show(r *gin.Engine) *gin.Engine {
	r.POST("/crud1", func(c *gin.Context) {

		price_val, err := strconv.ParseFloat(c.PostForm("price"), 32)
		if err != nil {
			panic(err)
		}

		quantity_val, err := strconv.Atoi(c.PostForm("quantity"))
		if err != nil {
			panic(err)
		}

		msg := db.Invoice{
			Fname:    c.PostForm("fname"),
			Lname:    c.PostForm("lname"),
			Product:  c.PostForm("product"),
			Price:    float32(price_val),
			Quantity: quantity_val,
			Category: c.PostForm("category"),
			Shipping: c.PostForm("shipping"),
		}

		c.JSON(http.StatusOK, msg)
	})
	return r
}

// renders form page to create an invoice
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
		c.HTML(http.StatusOK, "crud2.tmpl", gin.H{
			"title": "Crud2",
			"details": "Request the First Name, Last Name, Product," +
				"Price and Quantity of all the customers who's" +
				"invoice has a unit price over $13 ",
			"rslt": db.ReadOp(),
		})
	})
	return r
}

func main() {
	r := setRouter()
	r = readData(r)
	r = addInvoice(r)
	r = show(r)
	r.Run()
}
