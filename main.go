package main

import (
	"net/http"

	db "github.com/ScriptMang/conch/internal/bikeshop"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Conch Index-Page",
			"crud2": "Perform a Read Operation ",
		})
	})

	r.GET("/crud1", func(c *gin.Context) {
		c.HTML(http.StatusOK, "crud2.tmpl", gin.H{
			"title":   "Crud2",
			"details": "Request the First Name, Last Name, Product, Price and Quantity of all the customers who's invoice has a unit price over $13 ",
			"rslt":    db.ReadOp(),
		})
	})
	r.Run()

}
