package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

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

func show(r *gin.Engine) *gin.Engine {
	r.POST("/show", func(c *gin.Context) {

		// Convert the string values to float and int respectively
		// Then  create an Invoice Object based on the PostForm-Data

		c.Request.ParseForm()
		price_val, err := strconv.ParseFloat(c.PostForm("price"), 32)
		if err != nil {
			panic(err)
		}

		quantity_val, err := strconv.Atoi(c.PostForm("quantity"))
		if err != nil {
			panic(err)
		}

		toFile := db.Invoice{
			Fname:    c.PostForm("fname"),
			Lname:    c.PostForm("lname"),
			Product:  c.PostForm("product"),
			Price:    float32(price_val),
			Quantity: quantity_val,
			Category: c.PostForm("category"),
			Shipping: c.PostForm("shipping"),
		}

		// Create a New Encoder where you'll write your
		// json Encoding of toFile to the file stream tmpFile

		tmpFile, err := os.CreateTemp(os.TempDir(), "sample-data")

		if err != nil {
			panic(err)
		}

		defer os.Remove(tmpFile.Name())

		err = json.NewEncoder(tmpFile).Encode(toFile)
		if err != nil {
			panic(err)
		}

		err = tmpFile.Close()
		if err != nil {
			panic(err)
		}

		// Open that file Stream and
		// get a new Decoder stream to decode json as a string

		tmpFile2, err := os.Open(tmpFile.Name())
		if err != nil {
			panic(err)
		}

		var order db.Invoice

		err = json.NewDecoder(tmpFile2).Decode(&order)
		if err != nil {
			panic(err)
		}

		err = tmpFile2.Close()
		if err != nil {
			panic(err)
		}

		c.HTML(http.StatusOK, "show.tmpl", gin.H{
			"title": "Show Post Form Data",
			"rslt":  fmt.Sprintf("%+v\n", order),
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
	r = show(r)
	r.Run()
}
