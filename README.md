# bookstore-api

Demonstrate basic connection between mongo and golang.

Has following endpoint:

 - `/books` - GET - Get all books
 - `/books` - POST - Create new book, post a JSON body
 - `/books/isbn` - GET - Get the book
 - `/books/isbn` - PUT - Update the book, PUT a JSON body
 - `/books/isbn`- DELETE - Delete the book


## Development

As per dependency, _mgo needs 3.x version of mongod_ installed. This is likely to be updated as 4.x matures.

 - [ ] Add unittest.
