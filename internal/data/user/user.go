/*
Package user is a sample implementation of data repository.
It contains:
- Small "User" data model (id, name, surname);
- Repository abstraction and implementation with mongo-go-driver;
  - Get all, find one by id, delete all, delete one by id, insert new;
- Service abstraction and implementation (placeholder for business logic);
- HTTP endpoint "/users";
- Wrapper for mongo-go-driver package for mocking;
- Mocks and tests of repository implementation;
- Docker-compose files for launching with MongoDB;
 */
package user

//ID is string alias type for User id
type ID = string

//User represents a person
type User struct {
	ID      ID     `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
