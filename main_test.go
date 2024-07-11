package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	assert_v2 "github.com/go-playground/assert/v2"
)

func TestReadInvoice(t *testing.T) {

}

func TestPostInvoice(t *testing.T) {
	r := setRouter()
	r = show(r)
	w := httptest.NewRecorder()
	vals := url.Values{}
	vals.Add("fname", "johnny")
	vals.Add("lname", "TwoTap")
	vals.Add("product", "Peashooter")
	vals.Add("price", "20.00")
	vals.Add("quantity", "1")
	vals.Add("category", "Toy")
	vals.Add("shipping", "578 Bingus Ave, Moeberry OK 71203")

	sampleData := vals.Encode()
	fmt.Printf("Encoding: %v\n", sampleData)

	req, err := http.NewRequest("POST", "/crud1", strings.NewReader(sampleData))
	if err != nil {
		t.Fatalf("Error_v1:\n %v\n", err)
	}

	req.PostForm = vals
	r.ServeHTTP(w, req)

	expectedData := `{"fname":"johnny","lname":"TwoTap","product":"Peashooter","price":20,"quantity":1,"category":"Toy","shipping":"578 Bingus Ave, Moeberry OK 71203"}`

	rslt := w.Body.String()

	assert_v2.Equal(t, w.Code, 200)
	assert_v2.Equal(t, rslt, expectedData)
}
