# Conch

A locally ran program that uses gin to create api endpoints that perform CRUD Operations.
For this project I'm simulating a bikshop database that keeps a table of invoices.

### Prior to Running

* This program requires that both the go programming language
  and PostgreSQL be installed on your local machine.
  A program like Postman to send http requests to the database.

* Make sure to import the psql dump file into PostgreSQL.
  Use the following command 'psql --username=<username> Bikeshop <  <filename>.sql'.
  Note: replace placeholder info within the angle brackets with your own.

* Inside the bikeshop.go file replace the substring 'username' within the uri variable. 
  Provide the username you used to create the database in PostgreSQL instead.

### How to Run

The program is ran using the terminal. Typing 'go run .' or  by running 
its binary './conch' within in its directory, starts the program.
A prompt will show up  asking for outside internet connection, always deny it.


While the programs is still running, open Postman and send a request
using the following options below to run the crud operation.

Note: some CRUD operations will require you to pass JSON to the responsebody('Body' in PostMan) along with the request.
In that case, where you see <body> in 'CRUD Operations' replace it with the json object format listed below.
Also, do not send the id attribute as part of the json object, its created along with invoice.
For routes that end in '/:id' replace ':id' with an integer number.
To end the program in the terminal, type '^c'(ctrl-c).

### Expected Json Data
```
{
  "fname": string,
  "lname": string,
  "product": string,
  "price": float,
  "quantity": int,
  "category": string,
  "shipping": string
}
```

### CRUD Operations
* Add an invoice to the table<br>
   `POST` `localhost:8080/crud1/invoices/` `<body>`
* Read all the invoices from the table<br>
   `GET` `localhost:8080/crud2/invoices`
* Read an invoice based on their ID<br>
   `GET` `localhost:8080/crud2/invoice/:id`
* Edit an existing invoice<br>
   `PUT` `localhost:8080/crud3/invoice/:id` `<body>`
* Delete an existing invoice<br>
   `DELETE` `localhost:8080/crud4/invoice/:id`
