package storage

var _ = Entity(&User{})
var _ = Entities(Users{})

type User struct{ Email string }

func (u *User) Id() string { return u.Email }

type Users []*User
