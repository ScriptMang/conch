# Conch

A locally ran program that uses GIN Rest APIs to perform mock CRUD Operations.
For this project I'm simulating a bikshop database that keeps a table of invoices.

### Prior to Running

* This program requires that both the go programming language 
  and PostgreSQL be installed on your local machine.

* Make sure to import the psql dump file into PostgreSQL.
  Use the following command 'psql --username=<username> Bikeshop <  <filename>.sql'.
  Note: replace placeholder info within the angle brackets with your own.

* Inside the bikeshop.go file replace the substring 'username' within the uri variable. 
  Provide the username you used to create the database in PostgreSQL instead.

### How to Run

The program is ran using the terminal. Typing 'go run .' or  by running 
its binary './conch' within in its directory, starts the program.
A prompt will show up  asking for outside internet connection, always deny it.

Now open any browser and type 'localhost:8080' while the program is still running.
This address will take you the index page filled with a link for each CRUD Operation.
Clicking the link performs the crud operation unless it requires a html form to be sent first.
To end the program,  type '^c'(ctrl-c).

### THE CRUD Operations

* Create an invoice 
* Print the invoices table
* Edit an existing invoice
* Delete an existing invoice (WIP)
