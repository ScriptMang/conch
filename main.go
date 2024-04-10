package main

import (
	"fmt"
	"net/http"

	db "github.com/ScriptMang/conch/internal/bikeshop"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	fmt.Println("Hello World")
	db.New()
	r.Run()

}
