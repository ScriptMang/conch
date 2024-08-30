# Conch

A locally ran program that uses gin to create api endpoints that perform CRUD Operations.<br>
For this project I'm simulating a bikshop database that keeps a table of invoices.

### Prior to Running

* This program requires that both the go programming language<br>
  and PostgreSQL be installed on your local machine.<br>
  A program like Postman to send http requests to the database.

* Make sure to import the psql dump file into PostgreSQL.<br>
  Use the following command 'psql --username=<username> Bikeshop <  <filename>.sql'. <br>
  Note: replace placeholder info within the angle brackets with your own.

* Inside the bikeshop.go file replace the substring 'username' within the uri variable. <br>
  Provide the username you used to create the database in PostgreSQL instead.

### How to Run

The program is ran using the terminal. Typing `go run .`<br> 
or by running its binary `./conch` within in its directory, starts the program.<br>
A prompt will show up asking for outside internet connection, always deny it.

While the programs is still running, open Postman.<br>
Then, send a request using the following options below to run the crud operation.

#### Note About the CRUD Operations
some CRUD Operations will require you to pass JSON to the responsebody<br>
('Body' in PostMan) along with the request. In that case, where ever you see `<body>`<br> 
in 'CRUD Operations' replace it with the expected json format for an invoice listed below.<br>
Also, do not send the id attribute as part of the json object, its created along with the invoice.<br> 
For routes that end in `:id` replace it with an integer number.<br>
To end the program in the terminal, type `^c`(ctrl-c).

### Expected JSON Format for an Invoice
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
   `POST` `localhost:8080/invoices/` `<body>`
* Read all the invoices from the table<br>
   `GET` `localhost:8080/invoices`
* Read an invoice based on their ID<br>
   `GET` `localhost:8080/invoice/:id`
* Update an existing invoice<br>
   `PUT` `localhost:8080/invoice/:id` `<body>`
* Delete an existing invoice<br>
   `DELETE` `localhost:8080/invoice/:id`
