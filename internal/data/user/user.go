package user

//ID is id alias for User
type ID = string

//User represents a person
type User struct {
	ID      ID     `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
