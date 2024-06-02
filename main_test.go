package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ScriptMang/conch/internal/bikeshop"
	assert_v2 "github.com/go-playground/assert/v2"
)

func TestPostInvoice(t *testing.T) {
	r := setRouter()
	r = addInvoice(r)

	w := httptest.NewRecorder()

	sampleData := bikeshop.Invoice{
		Fname:    "Rocci",
		Lname:    "Marcia",
		Product:  "80lb barbell",
		Quantity: 1,
		Category: "exercise equipment",
		Shipping: "548 bukleet ave, NewPort MN 5510",
	}

	sampleJson, _ := json.Marshal(sampleData)
	req, _ := http.NewRequest("POST", "/crud1", strings.NewReader(string(sampleJson)))
	r.ServeHTTP(w, req)

	assert_v2.Equal(t, w.Code, 200)
	assert_v2.Equal(t, w.Body.String(), string(sampleJson))

}
