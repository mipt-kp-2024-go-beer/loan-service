# book-service
## Internal API (only for other microservices)
- Book info?
- Book get stock?

## Public API (may require auth)
- Book find: takes some criteria, returns a list of matching books.
- Book info: takes book id, returns book info.
- Book edit (requires permission): takes info.
- Book add (requires permission): takes info, returns id.
- Book delete (requires permission): takes id.
- Book get stock (requires permission): takes id, returns stock (including lent books).
- Book change stock (requires permission): takes id and delta, returns new stock.
- Some statistics?

# user-service
## Internal API (only for other microservices)
- Auth verify: takes token (+ permissions?), returns bool + user id (+ minimal info).

## Public API (may require auth)
- User login: takes credentials, returns a token.
- User create (might require permission): takes credentials and info, returns user id.
- User delete (requires permission / self?): takes user id. Might verify that the user has no unreturned books at the moment.
- User info (requires permission / self): takes user id, returns user info.
- User edit (requires permission / self?): takes info and updates it.
- User give permissions (requires permission): takes user id and new permissions.
- User revoke permissions (requires permission): takes user id and permissions to revoke.
- User check permissions? (requires permission?): takes user id and permissions, returns bool.

# loan-service
## Internal API (only for other microservices)
- User loan status: takes user id, returns whether the user has any unreturned books.

## Public API (may require auth)
- Book take (requires permission): takes book id (and optional user id if not for self).
- Book return (requires permission): takes book id (and optional user id if not for self).
- Book available (requires permission): takes book id, returns count left.
- Reservations list (requires permission): takes time, returns list of books taken at that point.
- Overdue list (requires permission): takes time margin, returns list of overdue books.
- Some statistics?
- Clean up database?
