package main

import (
	"net/http"

	db "github.com/ScriptMang/conch/internal/bikeshop"
	"github.com/gin-gonic/gin"
)

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

func addInvoice(r *gin.Engine) *gin.Engine {
	r.GET("/crud1", func(c *gin.Context) {
		c.HTML(http.StatusOK, "crud1.tmpl", gin.H{
			"title":   "Crud1",
			"details": "Add Invoice",
		})
	})
	return r
}

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
	r.Run()
}
