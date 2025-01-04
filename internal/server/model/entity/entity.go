// Package entity - entities
package entity

type User struct {
	ID       uint64
	Login    string
	PassHash string
}

type Secret struct {
	ID        uint64
	UserId    uint64
	UpdatedAt uint64
	Data      []byte
}
