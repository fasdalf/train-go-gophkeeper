// Package entity - entities
package entity

type User struct {
	ID       int64
	Login    string
	PassHash string
}

type Secret struct {
	ID        int64
	UserId    int64
	UpdatedAt int64
	IsDeleted bool
	Data      []byte
}
