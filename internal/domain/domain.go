package domain

import (
	"errors"
	"github.com/google/uuid"
)

type USERID = int64

type User struct {
	TgID           USERID
	ChosenFamilyID *uuid.UUID
	AccountName    string
	FullName       string
}

var (
	ErrUserExists = errors.New("user already exists")
)

type Family struct {
	ID   uuid.UUID
	Name string
}
