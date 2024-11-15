package fields

import (
	"strconv"
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

// checks the field to see if it exceeds or falls below a given char limit
// if it doesn't match the upper or lower limit an error message is added
// to the list of grammar errors
func isFieldTooLong(fieldName string, val *string, gramErr *GrammarError, minimum, maximum int) {
	fieldLen := len(*val)
	if fieldLen < minimum {
		gramErr.AddMsg(BadRequest, "Error: "+fieldName+" is too short, expected "+
			strconv.Itoa(minimum)+"-"+strconv.Itoa(maximum)+" chars")
	}
	if fieldLen > maximum {
		gramErr.AddMsg(BadRequest, "Error: "+fieldName+" is too long, expected "+
			strconv.Itoa(minimum)+"-"+strconv.Itoa(maximum)+" chars")
	}
}

// checks to see if there any capital letters in string val
// adds an new error to fieldErrs if none exist
func fieldHasNoCapLetters(val *string, fieldErr *GrammarError) {
	capLst := "ABCDEFGHIJKLMNOPQRYTUVWXYZ"
	if !strings.ContainsAny(*val, capLst) {
		fieldErr.AddMsg(BadRequest, "Error: Password must contain one or more capital letters")
	}
}

// checks to see if there are any digits in string val
// adds an new error to fieldErrs if none exist
func fieldHasNoNums(val *string, fieldErr *GrammarError) {
	nums := "0123456789"
	if !strings.ContainsAny(*val, nums) {
		fieldErr.AddMsg(BadRequest, "Error: Password must contain one or more digits")
	}
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

// swaps the orig values of an invoice for new ones
func CheckGrammarForPatch(val *string, fieldName string, orig string, fieldErr *GrammarError) {
	name := fieldName
	if *val == "" {
		// fmt.Printf("CheckGrammarForPatch: %s field value is blank\n", field.name)
		*val = orig // unique to patch requests
		// fmt.Println("CheckGrammarForPatch: Swap for orig.value: ", field.value)
	} else if *val != "" && name != "Address" &&
		name != "Product" && name != "Username" &&
		name != "Password" {
		fieldHasDigits(fieldName, val, fieldErr)
		fieldHasPunct(fieldName, val, fieldErr)
		fieldHasSymbols(fieldName, val, fieldErr)
	}

	if name == "Username" ||
		name == "Address" || name == "Product" {
		fieldHasPunct(fieldName, val, fieldErr)
		fieldHasSymbols(fieldName, val, fieldErr)
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
