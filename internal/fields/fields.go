package fields

import (
	"strings"
)

type GrammarError struct {
	ErrMsgs []string
}

var ErrorCode int // http-status code for errors
const BadRequest = 400
const ResourceNotFound = 404

// helper funct: takes a pointer to an InvoiceErorr, HttpStatusCode and a string msg
// as parameters and sets the values for the GrammarError struct.
// By default content-type is of type 'application/json'
func (fieldErr *GrammarError) AddMsg(statusCode int, str string) {
	ErrorCode = statusCode
	fieldErr.ErrMsgs = append(fieldErr.ErrMsgs, str)
}

// checks for empty text-fields in an invoice
// if there an error its added to an error slice
func isTextFieldEmpty(fieldName string, val *string, fieldErr *GrammarError) {
	if *val == "" {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't be empty")
	}
}

func fieldHasDigits(fieldName string, val *string, fieldErr *GrammarError) {
	digitFilter := "0123456789"
	if isTextInvalid(*val, digitFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any digits")
	}
}

func fieldHasPunct(fieldName string, val *string, fieldErr *GrammarError) {
	punctFilter := ".,?!'\"`:;"

	switch fieldName {
	case "Fname", "Lname":
		punctFilter = " .,?!'\"`:;"
	case "Product":
		punctFilter = "?!'\";"
	case "Category", "Address":
		punctFilter = ".?!'\"`:;"
	}

	if isTextInvalid(*val, punctFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any punctuation")
	}
}

func fieldHasSymbols(fieldName string, val *string, fieldErr *GrammarError) {
	symbolFilter := "~@#%$^|><&*()[]{}_-+=\\/"

	switch fieldName {
	case "Product":
		symbolFilter = "~#$*{}[]_\\+=><^"
	case "Category":
		symbolFilter = "~@#%$^|><*()[]{}_-+=\\/"
	case "Shipping":
		symbolFilter = "~@#&%$^|><*()[]{}_+=\\/"
	}

	// check for symbols: first-name, last-name, category, product
	if isTextInvalid(*val, symbolFilter) {
		fieldErr.AddMsg(BadRequest, "Error: "+fieldName+" can't have any Symbols")
	}
}

// checks a string field against an invalid char sequence
// if it returns a index then the text is invalid and it returns true
func isTextInvalid(val, charFilter string) bool {
	return strings.ContainsAny(val, charFilter)
}

// checks a field for punctuation, digits, and symbols
func CheckGrammar(fieldName string, val *string, fieldErr *GrammarError) {

	isTextFieldEmpty(fieldName, val, fieldErr)

	name := fieldName
	if *val != "" && name != "Address" &&
		name != "Product" && name != "Username" &&
		name != "Password" {
		fieldHasDigits(name, val, fieldErr)
		fieldHasPunct(name, val, fieldErr)
		fieldHasSymbols(name, val, fieldErr)
	}

	if name == "Username" ||
		name == "Address" || name == "Product" {
		fieldHasPunct(name, val, fieldErr)
		fieldHasSymbols(name, val, fieldErr)
	}

	if name == "Username" ||
		name == "Password" {
		isFieldTooLong(name, val, fieldErr, 8, 16)
	}

	if name == "Password" {
		fieldHasNoCapLetters(val, fieldErr)
		fieldHasNoNums(val, fieldErr)
	}
}

// func checkGrammarForPatch(field *textField, orig string, fieldErr *GrammarError) {
// 	name := field.name
// 	if *field.value == "" {
// 		//fmt.Printf("CheckGrammarForPatch: %s field value is blank\n", field.name)
// 		*field.value = orig // unique to patch requests
// 		//fmt.Println("CheckGrammarForPatch: Swap for orig.value: ", field.value)
// 	} else if *field.value != "" && name != "Shipping" && name != "Product" {
// 		fieldHasDigits(*field, fieldErr)
// 		fieldHasPunct(*field, fieldErr)
// 		fieldHasSymbols(*field, fieldErr)
// 	}

// 	if name == "Shipping" || name == "Product" {
// 		fieldHasPunct(*field, fieldErr)
// 		fieldHasSymbols(*field, fieldErr)
// 	}
// }
