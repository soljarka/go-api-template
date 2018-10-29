package user

//ID is string alias type for User id
type ID = string

//User represents a person
type User struct {
	ID      ID     `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}
