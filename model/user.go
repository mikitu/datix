package model

type User struct {
	Id uint64
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string	`json:"email"`
	Gender string	`json:"gender"`
	IpAddress string	`json:"ip_address"`
}

type Users []User
