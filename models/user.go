package models

type User struct {
	email    string
	password string
}

func CreateUser(email string, password string) User {
	return User{
		email:    email,
		password: password,
	}
}
func CreateUserPointer(email string, password string) *User {
	return &User{
		email:    email,
		password: password,
	}
}

//func InsertUser(user User) error {
//	return
//}

func (user *User) ToString() string {
	return user.email + " " + user.password
}
