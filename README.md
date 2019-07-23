# bookstore-api

Demonstrate basic connection between mongo and golang.

Has following endpoint:

 - `/books` - GET - Get all books
 - `/books` - POST - Create new book, post a JSON body
 - `/books/isbn` - GET - Get the book
 - `/books/isbn` - PUT - Update the book, PUT a JSON body
 - `/books/isbn`- DELETE - Delete the book

## Usage

### Add a new book

    curl -X POST -H "Content-Type: application/json" -d '{"isbn": "0134190440", "title": "The Go Programming Language", "authors": ["Alan A. A. Donovan", "Brian W. Kernighan"], "price": "$34.57"}' http://localhost:8080/books
    
### Get all books

    curl -H "Content-Type: application/json" http://localhost:8080/books

### Get a single book

Just pass the ISBN.

    curl -H "Content-Type: application/json" http://localhost:8080/books/0134190440

### Update a book

Pass JSON body to particular book endpoint

    curl -X PUT -H "Content-Type: application/json" -d '{"isbn": "0134190440", "title": "The Go Programming Language", "authors": ["Alan A. A. Donovan", "Brian W. Kernighan"], "price": "$20.00"}' http://localhost:8080/books/0134190440

### Delete a book

    curl -X DELETE -H "Content-Type: application/json" -d @body.json http://localhost:8080/books/0134190440

## Development

As per dependency, _mgo needs 3.x version of mongod_ installed. This is likely to be updated as 4.x matures.

 - [ ] Add unittest.
