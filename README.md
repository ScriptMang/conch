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
For routes that end in `:id` replace it with an integer number.<br>
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
  "product": string,
  "category": string,
  "price": float,
  "quantity": int
}
```


### Notes about the json objects

#### The id properties are never passed but created 
An Invoice structs id properties `:id` respectfully are assigned after creation.
They're incremented after insertion(pass or fail) by the database.

#### Passwords for Account structs get encrypted
After sending a Post request to create an account,
the password is encrypted and stored in the database.

### Basic Auth

Basic auth stands for basic authentication.
It requires a user to give their username
and password to authenticate. The client can
pass this info using curl by using the flag
to `authorization` and providing the login
credentials. The login credentials are combined 
together and  are encoded in base64.

You can use Postman click authorization, select basic-auth
and pass it the username and password have it encode it in
base64 and send it along with the `POST` request


### Tokens

Tokens are created after a user logs into their account. 
The token type being used is a bearer token.
The user must take their token and pass it across all
routes except account creation and logging in.
When a user logs out they delete their token.
This means to access any of their routes that needed
tokens before they need to login again and
generate a new token.




### CRUD Operations
* Add a user account to the table<br>
   `POST` `localhost:8080/users/` `<account>`
* Login into your account<br>
   `POST` `localhost:8080/login` `<basic-auth>`
* Log out of your account<br>
   `POST` `localhost:8080/logout` `<token>`
* Read all the usernames from the table<br>
   `GET` `localhost:8080/users`
* Read all the invoices from the table<br>
   `GET` `localhost:8080/invoices` `<token>`
* Read a specific user from the table<br>
   `GET` `localhost:8080/user` `<token>`
* Read all the invoices for a specific user<br>
   `GET` `localhost:8080/user/invoices` `<token>`
* Read an invoice for a specific user<br>
   `GET` `localhost:8080/invoice/:id` `<token>`
* Add an invoice to a specific user<br>
   `POST` `localhost:8080/invoices/` `<token>` `<invoice>`
* Update all the fields on an existing invoice<br>
   `PUT` `localhost:8080/invoice/:id`  `<token>` `<invoice>`
* Patch one or more fields on an existing invoice<br>
   `PATCH` `localhost:8080/invoice/:id` `<token>` `<invoice>`
* Delete an existing account<br>
   `DELETE` `localhost:8080/users` `<token>`
* Delete an existing invoice<br>
   `DELETE` `localhost:8080/invoice/:id` `<token>`
