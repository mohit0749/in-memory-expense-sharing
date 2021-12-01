package user

import (
	"sync/atomic"
)

var seqNo uint64

type User struct {
	id                   uint64
	Name, email, contact string
	lend, owes           []uint64
}

func NewUser(name, email, contact string) User {
	atomic.AddUint64(&seqNo, 1)
	return User{seqNo, name, email, contact, make([]uint64, 0),
		make([]uint64, 0)}
}

func (u User) GetId() uint64 {
	return u.id
}

func (u *User) Owe(expId uint64) {
	u.owes = append(u.owes, expId)
}

func (u *User) Lend(expId uint64) {
	u.lend = append(u.lend, expId)
}

func (u User) GetOwes() []uint64 {
	return u.owes
}

func (u User) GetLend() []uint64 {
	return u.lend
}
