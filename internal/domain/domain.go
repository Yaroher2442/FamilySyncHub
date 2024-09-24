package domain

import (
	"errors"
	"github.com/google/uuid"
)

type USERID = int64

//type UserState string
//
//const (
//	UserStateEmpty     UserState = "empty"
//	UserStateNewFamily UserState = "new_family"
//)

type User struct {
	TgID           USERID
	ChosenFamilyID *uuid.UUID
	AccountName    string
	FullName       string
	//State          UserState
}

var (
	ErrNoUser       = errors.New("user not found")
	ErrUserExists   = errors.New("user already exists")
	ErrFamilyExists = errors.New("family already exists")
)

type Family struct {
	ID   uuid.UUID
	Name string
}

var (
	ErrFamiliesEmpty = errors.New("user has no families")
)

type Category struct {
	ID       uuid.UUID
	Name     string
	FamilyID uuid.UUID
}

var (
	ErrDuplicateCategory = errors.New("category already exists")
)
