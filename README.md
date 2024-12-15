# Conch

A locally ran program that uses gin to create api endpoints that perform CRUD Operations.<br>
For this project I'm simulating a bikshop database that keeps a table of invoices.

### Prior to Running

* This program requires that both the go programming language<br>
  and PostgreSQL be installed on your local machine.<br>
  It also  requires a program like Postman to send http requests 
  to the database.

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
('Body' in PostMan) along with the request. In that case, where ever you see `<struct>`<br> 
in 'CRUD Operations' replace it with the expected json format for the object listed below.<br>
Also, do not send the id attribute as part of the json object, its created along with the invoice.<br> 
For routes that end in `:usr_id` or `:id` replace it with an integer number.<br>
To end the program in the terminal, type `^c`(ctrl-c).

### JSON Format for an Account
```
{
  "username": string,
  "fname": string,
  "lname": string,
  "address": string,
  "password": string
}
```

### JSON Format for a Usernames
```
{ 
  "id": int,
  "user_id": int,
  "username": string,
}
```

### JSON Format for a UserContacts
```
{ 
  "id": int,
  "user_id": int,
  "fname": string,
  "lname": string,
  "address": string
}
```

### JSON Format for an Invoice
```
{
  "id": int,
  "user_id": int,
  "product": string,
  "category": string,
  "price": float,
  "quantity": int
}
```

### JSON Format for Login-Creds
```
{
  "username": string,
  "pswd": []byte
}
```
### Notes about the json objects

#### Adding an account creates an entry in the Usernames table and the passwords table
Although you don't directly pass a json struct for a user
object in the request body, 
you still need to provide its id for other routes. i.e `:usr_id`.

#### The id properties are never passed but created 
User and Invoice structs' id properties `:usr_id` and `:id` respectfully are assigned after creation.
They're incremented after insertion(pass or fail) by the database.

#### Passwords for Account structs get encrypted
After sending a Post request to create an account,
the password is encrypted and stored in the database.


### CRUD Operations
* Add a user account to the table<br>
   `POST` `localhost:8080/users/` `<account>`
* Login into your account<br>
   `POST` `localhost:8080/user/login` `<login-creds>`
* Read all the usernames from the table<br>
   `GET` `localhost:8080/users`
* Read all the invoices from the table<br>
   `GET` `localhost:8080/invoices`
* Read a specific user from the table<br>
   `GET` `localhost:8080/user/:usr_id`
* Read all the invoices for a specific user<br>
   `GET` `localhost:8080/invoices`
* Read an invoice for a specific user<br>
   `GET` `localhost:8080/user/:usr_id/invoice/:id`
* Add an invoice to a specific user<br>
   `POST` `localhost:8080/invoices/`
* Update all the fields on an existing invoice<br>
   `PUT` `localhost:8080/user/:usr_id/invoice/:id` `<invoice>`
* Patch one or more fields on an existing invoice<br>
   `PATCH` `localhost:8080/user/:usr_id/invoice/:id` `<invoice>`
* Delete an existing account<br>
   `DELETE` `localhost:8080/users/` `<login-creds>`
* Delete an existing invoice<br>
   `DELETE` `localhost:8080/user/:usr_id/invoice/:id`
