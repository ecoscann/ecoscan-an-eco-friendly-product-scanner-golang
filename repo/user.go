package repo

type User struct{
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
	Points int `json:"points" db:"points"`
}