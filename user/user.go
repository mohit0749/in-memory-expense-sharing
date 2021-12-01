package user

import "sync/atomic"

var seqNo int64

type User struct {
	id                   int64
	Name, email, contact string
	lend, owes           []int64
}

func NewUser(name, email, contact string) User {
	atomic.AddInt64(&seqNo, 1)
	return User{seqNo, name, email, contact, make([]int64, 0), make([]int64, 0)}
}

func (u User) GetId() int64 {
	return u.id
}

func (u *User) Owe(expId int64) {
	u.owes = append(u.owes, expId)
}

func (u *User) Lend(expId int64) {
	u.lend = append(u.lend, expId)
}

func (u User) GetOwes() []int64 {
	return u.owes
}

func (u User) GetLend() []int64 {
	return u.lend
}
