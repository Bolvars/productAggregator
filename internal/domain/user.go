package domain

type User struct {
	id string
}

func NewUser(id string) *User {
	return &User{
		id: id,
	}
}

func (u *User) Id() string {
	return u.id
}
