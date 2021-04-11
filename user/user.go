package user

type User struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	City     string `json:"city" validate:"required"`
	Country  string `json:"country" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type Storage interface {
	CreateUser(newUser *User) error
	GetUserByEmail(email string) (*User, error)
}
