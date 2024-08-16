package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	db "github.com/ScriptMang/conch/internal/bikeshop"
	assert_v2 "github.com/go-playground/assert/v2"
)

func TestReadInvoice(t *testing.T) {
	r := setRouter()
	r = readData(r)
	w := httptest.NewRecorder()
	var s string
	req, err := http.NewRequest("GET", "/crud2", strings.NewReader(s))

	if err != nil {
		t.Fatalf("Error:\n %v\n", err)
	}

	r.ServeHTTP(w, req)
	expectedData := `{"fname": "Dante", "lname": "Ferges", "product": "Safety Goggles", "price": 15.99 , "quantity": 3, "category": "Safety Equipment", "shipping": "423 Elm St, Chicago IL 60629"},{"fname": "Michael", "lname": "Wither", "product": "Lubricant", "price": 11.99, "quantity": 1, "category": "Maintenance", "shipping": "230 Furginson Rd, Oklahoma OK 731303"},{"fname": Georgei, "lname": "Ventalin", "product": "Door Hinges", "price": 12.50, "quantity": 5, "category": "Home Improvement", "shipping": "495 Durvington Ave, Topeka KS 66603"},{"fname": "Edart", "lname": "Muskrat", "product": "Wrench", "price": 24.99, "quantity": 1, "category": "Plumbing", "shipping": "134 Pluton St, Boston MA 02108"},{"fname": "Abra", "lname": "Katern", "product": "DiscoBall", "price": 19.99, "quantity": 6, "category": "Party", "shipping": "829 Sherbet St, Portland ME 04102"},{"fname": Charles, "lname": "Tarly", "product": "Zombie book", "price": 14.99, "quantity": 2, "category": "Fiction", "shipping": "134 Pluton St, Boston MA 02108"} `

	rslt := w.Body.String()
	assert_v2.Equal(t, rslt, expectedData)

}

// test the creation of an invoice
// requires the table to empty and serial sequence to be reset to 1.
func TestPostInvoice(t *testing.T) {
	r := setRouter()
	r = addInvoice(r)
	w := httptest.NewRecorder()
	inv := db.Invoice{
		ID:       0,
		Fname:    "Johnny",
		Lname:    "TwoTap",
		Product:  "Peashooter",
		Price:    20.00,
		Quantity: 1,
		Category: "Toy",
		Shipping: "578 Bingus Ave, Moeberry OK 71203",
	}

	sampleData, err := json.Marshal(inv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to Encode  Object into json: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Encoding: %s\n", string(sampleData))

	req, err := http.NewRequest("POST", "/crud1", strings.NewReader(string(sampleData)))
	req.Header.Set("content-type", "application/json")

	if err != nil {
		t.Fatalf("Error_v1:\n %v\n", err)
	}

	r.ServeHTTP(w, req)

	expectedData := `{"id":1,"fname":"Johnny","lname":"TwoTap","product":"Peashooter","price":20,"quantity":1,"category":"Toy","shipping":"578 Bingus Ave, Moeberry OK 71203"}`

	rslt := w.Body.String()

	assert_v2.Equal(t, w.Code, http.StatusCreated)
	assert_v2.Equal(t, rslt, expectedData)
}
