// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package data

import ()

type Board struct {
	ID     int64
	Name   string
	Slug   string
	UserID int64
}

type Column struct {
	ID           int64
	Name         string
	ElementOrder int64
	BoardID      int64
}

type Item struct {
	ID           int64
	Name         string
	ElementOrder int64
	ColumnID     int64
}

type User struct {
	ID           int64
	Fullname     string
	Email        string
	PasswordHash string
}
